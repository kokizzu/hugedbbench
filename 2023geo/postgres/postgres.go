package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kokizzu/gotro/L"

	geo "hugedbbench/2023geo"
)

func main() {
	const connTpl = `postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d`
	connStr := fmt.Sprintf(connTpl,
		`root`,
		`password`,
		`127.0.0.1`,
		5432,
		`root`, // or postgres
		32,
	)

	ctx := context.Background()

	db, err := pgxpool.New(ctx, connStr)
	L.PanicIf(err, `pgxpool.Connect `+connStr)
	defer db.Close()

	const createTable = `CREATE TABLE IF NOT EXISTS points_sg (
		id BIGSERIAL PRIMARY KEY NOT NULL
		, lat FLOAT8 NOT NULL
		, long FLOAT8 NOT NULL
	)`

	_, err = db.Exec(ctx, createTable)
	L.PanicIf(err, `db.Exec `+createTable)

	const trunceTable = `TRUNCATE TABLE points_sg`
	_, err = db.Exec(ctx, trunceTable)
	L.PanicIf(err, `db.Exec `+trunceTable)

	const createFunc = `CREATE OR REPLACE FUNCTION distance(
    lat1 double precision,
    lon1 double precision,
    lat2 double precision,
    lon2 double precision)
  RETURNS double precision AS
$____$
DECLARE
    R integer = 6371e3; -- Meters
    rad double precision = 0.01745329252;

    φ1 double precision = lat1 * rad;
    φ2 double precision = lat2 * rad;
    Δφ double precision = (lat2-lat1) * rad;
    Δλ double precision = (lon2-lon1) * rad;

    a double precision = sin(Δφ/2) * sin(Δφ/2) + cos(φ1) * cos(φ2) * sin(Δλ/2) * sin(Δλ/2);
    c double precision = 2 * atan2(sqrt(a), sqrt(1-a));    
BEGIN                                                     
    RETURN R * c;        
END  
$____$
  LANGUAGE plpgsql VOLATILE
  COST 100;`
	_, err = db.Exec(ctx, createFunc)
	L.PanicIf(err, `db.Exec `+createFunc)

	geo.Insert100kPoints(func(lat float64, long float64, id uint64) error {
		_, err = db.Exec(ctx, `INSERT INTO points_sg (id, lat, long) VALUES ($1, $2, $3)`, id, lat, long)
		L.IsError(err, `db.Exec`)
		return err
	})

	geo.SearchRadius200k(func(lat float64, long float64, boxMeter float64, maxResult int64) (uint64, error) {
		delta := boxMeter / geo.DegToMeter / 2
		rows, err := db.Query(ctx, `
SELECT id, lat, long, distance($1,$2,lat,long) AS dist
FROM points_sg 
WHERE lat BETWEEN $1-$3 AND $2+$3 
  AND long BETWEEN $1-$3 AND $2+$3 
LIMIT $4
`,
			lat, long,
			delta,
			maxResult)
		if L.IsError(err, `db.Query`) {
			return 0, err
		}
		defer rows.Close()
		total := uint64(0)
		for rows.Next() {
			var id uint64
			var lat, long, dist float64
			err := rows.Scan(&id, &lat, &long, &dist)
			if L.IsError(err, `rows.Scan`) {
				return 0, err
			}
			total++
		}
		return total, nil
	})

	geo.MovingPoint(func(lat float64, long float64, id uint64) error {
		_, err := db.Exec(ctx, `
UPDATE points_sg SET lat = $1, long = $2
WHERE id = $3
`, lat, long, id)
		L.IsError(err, `db.Exec`)
		return err
	})
}
