package main

import (
	"fmt"
	"hugedbbench/2022fts/datasets"
	"time"

	"github.com/manticoresoftware/go-sdk/manticore"
)

const FtsName = `ManticoreSearch`
const GoRoutineCount = 1

// const indexName = `index1`
const RecordsPerInsert = 100000

//2580586
func main() {
	// udict := datasets.LoadUrbanDictionaryDatasets()
	cl := manticore.NewClient()
	cl.SetServer("127.0.0.1", 9308)
	_, err := cl.Open()
	if err != nil {
		panic(err)
	}
	_, err = cl.Sphinxql(`DROP TABLE urbandict`)
	if err != nil {
		panic(err)
	}
	// fmt.Println(res)

	//id already exists by default
	// res, err = cl.Sphinxql(`CREATE TABLE massive (_id text,locale text, partition text, scenario text, intent text,utt text, annot_utt text, worker_id text)`)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)
	// type MassiveDatasets struct {
	// 	Id        string `json:"id"`
	// 	Locale    string `json:"locale"`
	// 	Partition string `json:"partition"`
	// 	Scenario  string `json:"scenario"`
	// 	Intent    string `json:"intent"`
	// 	Utt       string `json:"utt"`
	// 	AnnotUtt  string `json:"annot_utt"`
	// 	WorkerId  string `json:"worker_id"`
	// }

	// type UrbanDictionary struct {
	// 	WordId     string `json:"word_id"`
	// 	Word       string `json:"word"`
	// 	UpVotes    string `json:"up_votes"`
	// 	DownVotes  string `json:"down_votes"`
	// 	Author     string `json:"author"`
	// 	Definition string `json:"definition"`
	// }
	res, err := cl.Sphinxql(`CREATE TABLE urbandict (word_id text,word text, up_votes text, down_votes text, author text,definition text)`)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	createSqlValue := func(s, suffix string) string {
		return `'` + s + `'` + suffix
	}

	//create a fresh table

	reader := datasets.UrbanDictionaryReader{SkipHeader: true}
	insertDuration := []int{}
	for {
		datasets, _, err := reader.ReadNextNLines(100)

		q := `INSERT INTO urbandict (word_id, word, up_votes, down_votes, author,definition) Values `
		start := time.Now()
		for _, v := range datasets {
			q += `(` + createSqlValue(v.WordId, ` ,`) + createSqlValue(v.Word, ` ,`) + createSqlValue(v.UpVotes, ` ,`) + createSqlValue(v.DownVotes, ` ,`) + createSqlValue(v.Author, ` ,`) + createSqlValue(v.Definition, ``) + ` ),`
		}
		res, _ = cl.Sphinxql(q[:len(q)-1])
		insertDuration = append(insertDuration, int(time.Since(start)))
		if res[0].ErrorCode > 0 {
			fmt.Println(res)

			return
		}
		if err != nil {
			break
		}
		// fmt.Println(res, err)
		// res, err = cl.Sphinxql(q)
		// fmt.Println(res, err)
	}
	fmt.Println(FtsName+` BulkInsert 100 `, time.Duration(average(insertDuration)))
	fmt.Println(FtsName+` TotalInsert `, time.Duration(total(insertDuration)))

	// res, err = cl.Sphinxql(`replace into testrt values(2,'another subject', 'more content', 15)`)
	// fmt.Println(res, err)
	// res, err = cl.Sphinxql(`replace into testrt values(5,'again subject', 'one more content', 10)`)
	// fmt.Println(res, err)
	// start := time.Now()
	s := manticore.NewSearch(`anime`, `urbandict`, ``)
	s.Limit = 1000

	res2, err2 := cl.RunQuery(s)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	// fmt.Println(FtsName+`search`, time.Since(start))
	fmt.Println(FtsName+`search`, res2.QueryTime)
	fmt.Println(res2.Total)
	fmt.Println(res2.TotalFound)
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
