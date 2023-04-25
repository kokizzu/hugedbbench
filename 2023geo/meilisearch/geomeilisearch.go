package main

import (
	"fmt"

	"github.com/kokizzu/gotro/L"
	"github.com/meilisearch/meilisearch-go"

	geo "hugedbbench/2023geo"
)

func main() {

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://127.0.0.1:7700",
		APIKey: "123",
	})

	index := client.Index("points_sg")

	//L.Print(`DeleteAllDocuments`)
	//task, err := index.DeleteAllDocuments()
	//L.PanicIf(err, `index.DeleteAllDocuments`)
	//_, err = client.WaitForTask(task.TaskUID)
	//L.PanicIf(err, `client.WaitForTask DeleteAllDocuments`)
	//
	//L.Print(`UpdateSettings`)
	//task, err = index.UpdateSettings(&meilisearch.Settings{
	//	FilterableAttributes: []string{`_geo`},
	//	SortableAttributes:   []string{`_geo`},
	//})
	//L.PanicIf(err, `index.UpdateSettings`)
	//_, err = client.WaitForTask(task.TaskUID)
	//L.PanicIf(err, `client.WaitForTask UpdateSettings`)
	//
	//L.Print(`StartBenchmark`)
	//
	//// not optimal for non-batch request, just like clickhouse
	//geo.Insert100kPoints(func(lat, long float64, id uint64) error {
	//	task, err := index.AddDocuments([]map[string]interface{}{
	//		{"id": id, "_geo": []any{lat, long}},
	//	})
	//	L.IsError(err, `index.AddDocuments`)
	//	if err == nil {
	//		_, err := client.WaitForTask(task.TaskUID)
	//		L.IsError(err, `client.WaitForTask AddDocuments`)
	//		return err
	//	}
	//	return err
	//})

	geo.SearchRadius200k(func(lat, long, boxMeter float64, maxResult int64) (uint64, error) {
		// _geoBoundingBox: [minLat, minLong, maxLat, maxLong] --> didn't work
		res, err := index.Search("", &meilisearch.SearchRequest{
			Filter: fmt.Sprintf("_geoRadius(%f, %f, %.0f)", lat, long, boxMeter),
			Sort:   []string{fmt.Sprintf(`_geoPoint(%f, %f):asc`, lat, long)},
			Limit:  maxResult,
		})
		if L.IsError(err, `index.Search`) {
			return 0, err
		}
		// TODO: need to calculate distance manually (~10% overhead)
		return uint64(len(res.Hits)), nil
	})

	geo.MovingPoint(func(lat, long float64, id uint64) error {
		_, err := index.UpdateDocuments([]map[string]interface{}{
			{"id": id, "_geo": []any{lat, long}},
		})
		L.IsError(err, `index.UpdateDocuments`)
		return err
	})
}
