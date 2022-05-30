package main

import (
	"fmt"
	"hugedbbench/2022fts/datasets"

	"github.com/expectedsh/go-sonic/sonic"
)

func main() {
	//not sure about this one
	udict := datasets.LoadUrbanDictionaryDatasets()
	ingester, err := sonic.NewIngester("localhost", 1491, "SecretPassword")
	if err != nil {
		panic(err)
	}

	bulk := make([]sonic.IngestBulkRecord, len(udict))
	for i, v := range udict {
		bulk[i] = sonic.IngestBulkRecord{Object: v.WordId, Text: v.Definition}
	}
	_ = ingester.BulkPush("myIndex1", "general", 3, bulk, sonic.LangEng)
	search, err := sonic.NewSearch("localhost", 1491, "SecretPassword")
	if err != nil {
		panic(err)
	}

	results, _ := search.Query("myIndex1", "general", "anime", 10, 0, sonic.LangEng)

	fmt.Println(results)
}
