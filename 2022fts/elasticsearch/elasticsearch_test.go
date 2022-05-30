package main

import (
	"context"
	"fmt"
	"hugedbbench/2022fts/datasets"
	"testing"
	"time"

	"github.com/kokizzu/gotro/D/Es"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/olivere/elastic/v7"
)

func TestElasticSearch(t *testing.T) {
	// urbanDict := datasets.LoadUrbanDictionaryDatasets()[:50000]
	// if len(urbanDict) == 0 {
	// 	t.Error(`empty datasets`)
	// 	return
	// }
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(`http://127.0.0.1:9200`),
		elastic.SetSniff(false),
		// elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
		// elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
	}
	esClient, err := elastic.NewClient(opts...)
	L.PanicIf(err, `elastic.NewClient`)
	adapter := Es.Adapter{
		Client: esClient,
		Reconnect: func() *elastic.Client {
			esClient, err := elastic.NewClient(opts...)
			L.IsError(err, `elastic.NewClient %v`, opts)
			return esClient
		},
	}
	ctx := context.Background()
	// adapter.CreateIndex(IndexName)

	// ch := make(chan datasets.UrbanDictionary, RecordsPerGoroutine)
	// insert
	// wg := sync.WaitGroup{}
	// ch := make(chan datasets.UrbanDictionary)
	// go datasets.LoadUrbanDictionaryDatasetsChan(ch)
	adapter.DeleteIndex(IndexName).Do(ctx)
	reader := datasets.UrbanDictionaryReader{SkipHeader: true}
	bulkReq := adapter.Bulk()
	start := time.Now()
	t.Run(`insert`, func(t *testing.T) {
		for {
			datasets, _, err := reader.ReadNextNLines(100)
			if err != nil {
				break
			}
			for i := 0; i < len(datasets); i++ {
				bulkReq.Add(elastic.NewBulkIndexRequest().Index(IndexName).Doc(datasets[i]))
			}
			bulkReq.Refresh(`true`).Do(ctx)
			adapter.Reindex().Do(ctx)
		}
		// for i := 0; i < GoRoutineCount; i++ {
		// 	bulkReq := adapter.Bulk()
		// 	wg.Add(1)
		// 	go func(br *elastic.BulkService) {
		// 		for j := 1; ; j++ {
		// 			data, ok := <-ch
		// 			if !ok {
		// 				break
		// 			}
		// 			br.Add(elastic.NewBulkIndexRequest().Index(IndexName).Doc(data))
		// 			if j%RecordsPerInsert == 0 {
		// 				br.Refresh(`wait_for`).Do(ctx)
		// 			}
		// 		}

		// 		br.Refresh(`wait_for`).Do(ctx)
		// 		wg.Done()
		// 	}(bulkReq)
		// }
		// wg.Wait()
	})
	// adapter.Refresh(IndexName).Do(ctx)
	L.PanicIf(err, `bulkReq.Do`)
	fmt.Println(FtsName+` BulkInsert `, time.Since(start))
	start = time.Now()
	res := []string{}
	adapter.QueryRaw(IndexName, M.SX{
		`query`: M.SX{
			`bool`: M.SX{
				`should`: []interface{}{
					// M.SX{`match`: M.SX{`word`: `wood`}},
					M.SX{`match`: M.SX{`definition`: `anime`}},
					// M.SX{`match`: M.SX{`partition`: `test`}},
				},
			},
		},
	}, func(id string, rawJson []byte) (exitEarly bool) {
		res = append(res, id)
		return false
	})
	fmt.Println(FtsName+` Search`, time.Since(start))
	fmt.Println(`total: `, len(res))
}
