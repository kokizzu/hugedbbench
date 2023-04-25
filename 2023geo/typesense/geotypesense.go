package main

import (
	"fmt"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/typesense/typesense-go/typesense"
	_ "github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"

	geo "hugedbbench/2023geo"
)

func main() {
	client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("123"),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)

	const collectionName = "points_sg"
	// drop schema
	_, err := client.Collection(collectionName).Delete()
	if err != nil {
		if S.Contains(err.Error(), `No collection with name`) {
			err = nil
		}
		L.PanicIf(err, `client.Collection(collectionName).Delete()`)
	}

	// create schema
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{
				Name: "id",
				Type: "string", // cannot be int64
			},
			{
				Name: "coord",
				Type: "geopoint",
			},
		},
	}
	_, err = client.Collections().Create(schema)
	L.PanicIf(err, `client.Collections().Create(schema)`)

	doc := client.Collection(collectionName).Documents()

	geo.Insert100kPoints(func(lat, long float64, id uint64) error {
		_, err := doc.Create(M.SX{
			`id`:    fmt.Sprint(id),
			`coord`: []any{lat, long},
		})
		L.IsError(err, `doc.Create`)
		return err
	})

	geo.SearchRadius200k(func(lat, long, boxMeter float64, maxResult int64) (uint64, error) {
		sp := &api.SearchCollectionParams{
			Q: "*",
			// cannot be meter, must be km or mi
			FilterBy: pointer.String(fmt.Sprintf("coord:(%f,%f, %.2f km)", lat, long, boxMeter/1000)),
			SortBy:   pointer.String(fmt.Sprintf("coord(%f,%f):asc", lat, long)),
			PerPage:  pointer.Int(250), // max 250
		}
		res, err := doc.Search(sp)
		if L.IsError(err, `doc.Search`) {
			return 0, err
		}
		return uint64(len(*res.Hits)), nil
	})

	geo.MovingPoint(func(lat, long float64, id uint64) error {
		_, err := doc.Upsert(M.SX{
			`id`:    fmt.Sprint(id),
			`coord`: []any{lat, long},
		})
		L.IsError(err, `doc.Upsert`)
		return err
	})
}
