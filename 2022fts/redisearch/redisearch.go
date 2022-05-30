package main

import (
	"fmt"
	"hugedbbench/2022fts/datasets"
	"log"
	"strconv"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
)

const FtsName = `Redisearch`

func main() {

	// Create a client. By default a client is schemaless
	// unless a schema is provided when creating the index
	c := redisearch.NewClient("localhost:6379", "myIndex")
	// udict := datasets.LoadUrbanDictionaryDatasets()[:100000]

	// Create a schema

	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("word_id")).
		AddField(redisearch.NewTextFieldOptions("word", redisearch.TextFieldOptions{Weight: 5.0, Sortable: true})).
		AddField(redisearch.NewTextField("up_votes")).
		AddField(redisearch.NewTextField("down_votes")).
		AddField(redisearch.NewTextField(`author`)).
		AddField(redisearch.NewTextFieldOptions(`definition`, redisearch.TextFieldOptions{Weight: 5.0, Sortable: true}))
	// Drop an existing index. If the index does not exist an error is returned

	c.Drop()
	// Create the index with the given schema
	if err := c.CreateIndex(sc); err != nil {
		log.Fatal(err)
	}
	// Create a document with an id and given score
	reader := datasets.UrbanDictionaryReader{SkipHeader: true}

	idx := 1
	insertDuration := []int{}
	for {
		datasets, _, err := reader.ReadNextNLines(100)
		docs := make([]redisearch.Document, len(datasets))
		start := time.Now()
		for i, v := range datasets {
			doc := redisearch.NewDocument(strconv.Itoa(idx), 1.0)
			idx++
			doc.Set(sc.Fields[0].Name, v.WordId)
			doc.Set(sc.Fields[1].Name, v.Word)
			doc.Set(sc.Fields[2].Name, v.UpVotes)
			doc.Set(sc.Fields[3].Name, v.DownVotes)
			doc.Set(sc.Fields[4].Name, v.Author)
			doc.Set(sc.Fields[5].Name, v.Definition)
			docs[i] = doc
		}
		if err := c.Index(docs...); err != nil {
			log.Fatal(err)
		}
		insertDuration = append(insertDuration, int(time.Since(start)))
		if err != nil {
			break
		}
	}
	fmt.Println(FtsName+` BulkInsert 100 `, time.Duration(average(insertDuration)))
	fmt.Println(FtsName+` TotalInsert `, time.Duration(total(insertDuration)))
	// Searching with limit and sorting
	start := time.Now()
	docs, total, err := c.Search(redisearch.NewQuery(`anime`).
		AddFilter(redisearch.Filter{Field: sc.Fields[5].Name}).
		Limit(0, 1000))

	fmt.Println(FtsName+` search `, time.Since(start))
	if err != nil {
		panic(err)
	}
	fmt.Println(len(docs), total, err)
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
