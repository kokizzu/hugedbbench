package main

import (
	"fmt"
	"hugedbbench/2022fts/datasets"
	"log"
	"time"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

const FtsName = `typesense`

func main() {
	client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("local-typesense-api-key"),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)
	//Create a collection

	yes := true
	// sortField := "num_employees"
	indexName := `UrbanDictionary`
	// udict := datasets.LoadUrbanDictionaryDatasets()[:100000]
	schema := &api.CollectionSchema{
		Name: indexName,
		Fields: []api.Field{
			{
				Name: "word_id",
				Type: "string",
			},
			{
				Name:  "word",
				Type:  "string",
				Index: &yes,
			},
			{
				Name: "up_votes",
				Type: "string",
			},
			{
				Name: "down_votes",
				Type: "string",
			},
			{
				Name:  "author",
				Type:  "string",
				Index: &yes,
			},
			{
				Name:  "definition",
				Type:  "string",
				Facet: &yes,
				Index: &yes,
			},
		},
		// DefaultSortingField: &sortField,
	}
	client.Collection(indexName).Delete()
	client.Collections().Create(schema)

	//Index a document

	// document := struct {
	// 	ID           string `json:"id"`
	// 	CompanyName  string `json:"company_name"`
	// 	NumEmployees int    `json:"num_employees"`
	// 	Country      string `json:"country"`
	// }{
	// 	ID:           "123",
	// 	CompanyName:  "Stark Industries",
	// 	NumEmployees: 5215,
	// 	Country:      "USA",
	// }
	reader := datasets.UrbanDictionaryReader{SkipHeader: true}
	insertDuration := []int{}
	for {
		datasets, _, err := reader.ReadNextNLines(100)
		start := time.Now()
		toArrInterface := make([]interface{}, len(datasets))
		for idx, v := range datasets {
			toArrInterface[idx] = v
		}
		actualSize := len(datasets)
		//no bulk insert instead use import. see https://github.com/typesense/typesense/issues/35
		_, errx := client.Collection(indexName).Documents().Import(toArrInterface, &api.ImportDocumentsParams{BatchSize: &actualSize})
		insertDuration = append(insertDuration, int(time.Since(start)))
		if errx != nil {
			log.Println(errx)
		}
		if err != nil {
			break
		}
	}
	fmt.Println(FtsName+` BulkInsert 100 `, time.Duration(average(insertDuration)))
	fmt.Println(FtsName+` TotalInsert `, time.Duration(total(insertDuration)))

	limit := 1000
	searchParameters := &api.SearchCollectionParams{
		Q:          "anime",
		QueryBy:    "definition",
		GroupLimit: &limit,
	}

	res, err := client.Collection(indexName).Documents().Search(searchParameters)
	if err != nil {
		panic(err)
	}
	fmt.Println(*res.Found)
	fmt.Println(FtsName+` search `, *res.SearchTimeMs, ` ms`)
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
