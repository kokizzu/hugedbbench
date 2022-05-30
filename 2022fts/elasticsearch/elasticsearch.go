package main

import (
	"context"
	"fmt"
	"hugedbbench/2022fts/datasets"
	"log"
	"os"
	"time"

	"github.com/kokizzu/gotro/D/Es"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"

	"github.com/olivere/elastic/v7"
)

const FtsName = `ElasticSearch`
const GoRoutineCount = 1
const IndexName = `index1`
const RecordsPerInsert = 100000

//
func main() {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(`http://127.0.0.1:9200`),
		elastic.SetSniff(false),
	}
	opts = append(opts,
		elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
	)
	// urbanDict := datasets.LoadUrbanDictionaryDatasets()[:50000]
	// if len(urbanDict) == 0 {
	// 	t.Error(`empty datasets`)
	// 	return
	// }
	// opts := []elastic.ClientOptionFunc{
	// 	elastic.SetURL(`http://127.0.0.1:9200`),
	// 	elastic.SetSniff(false),
	// 	// elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
	// 	// elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
	// }
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
	adapter.DeleteIndex(IndexName).Do(ctx)
	reader := datasets.UrbanDictionaryReader{SkipHeader: true}
	bulkReq := adapter.Bulk()

	insertDuration := []int{}
	for {
		datasets, _, err := reader.ReadNextNLines(100)
		start := time.Now()
		for _, v := range datasets {
			bulkReq.Add(elastic.NewBulkIndexRequest().Index(IndexName).Doc(v))
		}

		bulkReq.Refresh(`true`).Do(ctx)
		adapter.Reindex().Do(ctx)
		insertDuration = append(insertDuration, int(time.Since(start)))
		if err != nil {
			break
		}
	}

	L.PanicIf(err, `bulkReq.Do`)
	fmt.Println(FtsName+` BulkInsert 100 `, time.Duration(average(insertDuration)))
	fmt.Println(FtsName+` TotalInsert `, time.Duration(total(insertDuration)))
	start := time.Now()
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

func total(x []int) int {
	i := 0
	for _, v := range x {
		i += v
	}
	return i
}

func average(x []int) int {
	i := 0
	for _, v := range x {
		i += v
	}
	return i / len(x)
}
