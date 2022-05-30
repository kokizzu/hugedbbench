package main

import (
	"fmt"
	"hugedbbench/2022fts/datasets"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

const FtsName = `Meilisearch`

func main() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    "http://127.0.0.1:7720",
		APIKey:  "test_api_key",
		Timeout: 600 * time.Second,
	})
	// An index is where the documents are stored.
	index := client.Index("urbandict")
	index.DeleteAllDocuments()
	reader := datasets.UrbanDictionaryReader{SkipHeader: true}

	insertDuration := []int{}

	for {
		datasets, _, err := reader.ReadNextNLines(100)
		start := time.Now()
		//i hope this working as expected
		t, _ := index.AddDocuments(datasets)
		//around 3.5min each https://github.com/meilisearch/MeiliSearch/issues/1098
		index.WaitForTask(&meilisearch.Task{UID: t.UID, Status: meilisearch.TaskStatusSucceeded})
		insertDuration = append(insertDuration, int(time.Since(start)))
		if err != nil {
			break
		}
	}
	req := []string{`definition`}
	index.UpdateSearchableAttributes(&req)
	fmt.Println(FtsName+` BulkInsert 100 `, time.Duration(average(insertDuration)))
	fmt.Println(FtsName+` TotalInsert `, time.Duration(total(insertDuration)))

	// start := time.Now()
	sr, err := index.Search(`anime`, &meilisearch.SearchRequest{Limit: 1000})
	if err != nil {
		panic(err)
	}
	fmt.Println(`found`, len(sr.Hits))

	//fmt.Println(FtsName+`search `,time.Since(start) )
	//i think using ProcessingTimeMs directly is better than time
	fmt.Println(FtsName+`search `, sr.ProcessingTimeMs, ` ms`)
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
