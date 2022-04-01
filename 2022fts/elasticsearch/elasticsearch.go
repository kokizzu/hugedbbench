package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kokizzu/gotro/D/Es"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/olivere/elastic/v7"
	_ "github.com/olivere/elastic/v7"
)

func main() {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(`http://127.0.0.1:9200`),
		elastic.SetSniff(false),
	}
	opts = append(opts,
		elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
	)
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

	const indexName = `index1`
	ctx := context.Background()

	// insert
	bulkReq := adapter.Bulk()
	bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`1`).Doc(M.SX{
		`type1`: `B`,
		`type2`: `Y`,
	}))
	bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`2`).Doc(M.SX{
		`type1`: `C`,
		`type2`: `Z`,
	}))
	bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`3`).Doc(M.SX{
		`type1`: `B`,
		`type2`: `Z`,
	}))
	bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`4`).Doc(M.SX{
		`type1`: `C`,
		`type2`: `X`,
	}))
	bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`5`).Doc(M.SX{
		`type1`: `A`,
		`type2`: `X`,
	}))
	_, err = bulkReq.Refresh("true").Do(ctx)
	L.PanicIf(err, `bulkReq.Do`)

	// search
	res := []string{}
	adapter.QueryRaw(indexName, M.SX{
		`query`: M.SX{
			`bool`: M.SX{
				`should`: []interface{}{
					M.SX{`match`: M.SX{`type2`: `Z`}},
					M.SX{`match`: M.SX{`type1`: `B`}},
					M.SX{`match`: M.SX{`type1`: `C`}},
				},
			},
		},
	}, func(id string, rawJson []byte) (exitEarly bool) {
		res = append(res, id)
		return false
	})
	fmt.Println(res)
}
