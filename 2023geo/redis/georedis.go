package main

import (
	"context"
	"fmt"

	"github.com/kokizzu/gotro/L"
	"github.com/rueian/rueidis"

	geo "hugedbbench/2023geo"
)

func main() {
	cli, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{`127.0.0.1:6379`},
		Password:    `kl234j23095125125125`,
	})
	L.PanicIf(err, `rueidis.NewClient`)
	defer cli.Close()

	ctx := context.Background()

	const key = `SG` // per city

	// clear databae
	b := cli.B().Del().Key(key).Build()
	resp := cli.Do(ctx, b)
	L.IsError(resp.Error(), `failed to DEL`, key)

	geo.Insert100kPoints(func(lat, long float64, id uint64) error {
		b := cli.B().Geoadd().Key(key).
			LongitudeLatitudeMember().LongitudeLatitudeMember(long, lat, fmt.Sprint(id)).Build()
		resp := cli.Do(ctx, b)
		err := resp.Error()
		L.IsError(err, `GEOADD`)
		return err
	})

	geo.SearchRadius200k(func(lat, long, boxMeter float64, maxResult int64) (uint64, error) {
		b := cli.B().Geosearch().Key(key).
			Fromlonlat(long, lat).
			Bybox(boxMeter).Height(boxMeter).M().
			Asc().
			Count(maxResult).
			Withcoord().
			Withdist().
			Build()
		resp := cli.Do(ctx, b)
		err := resp.Error()
		if L.IsError(err, `GEOSEARCH`) {
			return 0, err
		}
		rows, err := resp.ToArray()
		if L.IsError(err, `resp.ToArray`) {
			return 0, err
		}
		for _, row := range rows {
			col, err := row.ToArray()
			if L.IsError(err, `row.ToArray`) {
				return 0, err
			}
			_, err = col[0].ToString() // id
			if L.IsError(err, `col[0].ToString`) {
				return 0, err
			}
			_, _ = col[1].AsFloat64() // distance
			if L.IsError(err, `col[1].ToFloat64`) {
				return 0, err
			}
			coord, _ := col[2].ToArray()
			_, err = coord[0].AsFloat64() // long
			if L.IsError(err, `coord[0].ToFloat64`) {
				return 0, err
			}
			_, _ = coord[1].AsFloat64() // lat
			if L.IsError(err, `coord[1].ToFloat64`) {
				return 0, err
			}
		}
		return uint64(len(rows)), err
	})

	geo.MovingPoint(func(lat, long float64, id uint64) error {
		b := cli.B().Geoadd().Key(key).Xx().
			LongitudeLatitudeMember().LongitudeLatitudeMember(long, lat, fmt.Sprint(id)).Build()
		resp := cli.Do(ctx, b)
		err := resp.Error()
		L.IsError(err, `GEOADD XX`)
		return err
	})
}
