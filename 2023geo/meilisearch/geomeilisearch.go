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
	//called := atomic.Uint32{}
	//geo.Insert100kPoints(func(lat, long float64, id uint64) error {
	//	called.Add(1)
	//	task, err := index.AddDocuments([]map[string]interface{}{
	//		{"id": id, "_geo": []any{lat, long}},
	//	})
	//	L.IsError(err, `index.AddDocuments`)
	//	// wait only last one to make it a bit faster for Meilisearch
	//	if err == nil && called.Load() == 100_000 {
	//		_, err := client.WaitForTask(task.TaskUID)
	//		L.IsError(err, `client.WaitForTask AddDocuments`)
	//		return err
	//	}
	//	return err
	//})

	geo.SearchRadius200k(func(lat, long, boxMeter float64, maxResult int64) (uint64, error) {
		// _geoBoundingBox: --> didn't work
		delta := boxMeter / geo.DegToMeter / 2
		lat1 := lat - delta
		lat2 := lat + delta
		long1 := long - delta
		long2 := long + delta
		//fmt.Sprintf("_geoRadius(%f, %f, %.0f)", lat, long, boxMeter),
		res, err := index.Search("", &meilisearch.SearchRequest{
			Filter: fmt.Sprintf("_geoBoundingBox([%f, %f], [%f, %f])", lat2, long2, lat1, long1),
			Sort:   []string{fmt.Sprintf(`_geoPoint(%f, %f):asc`, lat, long)},
			Limit:  maxResult,
		})
		if L.IsError(err, `index.Search`) {
			return 0, err
		}
		rows := make([]any, 0, len(res.Hits))
		for _, row := range res.Hits {
			rows = append(rows, []any{row}) // TODO: get id, lat, long, _geoDistance (remove unecessary fields)
		}
		return uint64(len(rows)), nil
	})

	geo.MovingPoint(func(lat, long float64, id uint64) error {
		_, err := index.UpdateDocuments([]map[string]interface{}{
			{"id": id, "_geo": []any{lat, long}},
		})
		L.IsError(err, `index.UpdateDocuments`)
		return err
	})
}
