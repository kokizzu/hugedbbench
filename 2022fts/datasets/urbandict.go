package datasets

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/kokizzu/gotro/S"
)

//datasets https://www.kaggle.com/datasets/therohk/urban-dictionary-words-dataset
//go:generate gomodifytags -file urbandict.go -struct UrbanDictionaryDatasets -add-tags json -w
type UrbanDictionaryDatasets struct {
	WordId     string `json:"word_id"`
	Word       string `json:"word"`
	UpVotes    string `json:"up_votes"`
	DownVotes  string `json:"down_votes"`
	Author     string `json:"author"`
	Definition string `json:"definition"`
}

type UrbanDictionaryReader struct {
	SkipHeader  bool
	initialized bool

	index int64

	ErrorCount int64
	reader     *csv.Reader
	mutex      sync.Mutex
}

func (d *UrbanDictionaryReader) init() {
	d.mutex.Lock()
	if d.initialized {
		d.mutex.Unlock()
		return
	}
	_, fileDir, _, _ := runtime.Caller(0)
	//getting dir and concat with datasets
	f, err := os.OpenFile(fileDir[:len(fileDir)-len(`urbandict.go`)]+`urbandict-word-defs.csv`, os.O_RDONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	rd := csv.NewReader(f)
	if d.SkipHeader {
		rd.Read()
	}
	d.reader = rd
	d.initialized = true
	d.mutex.Unlock()
}

/*
ReadNextNLines return the datasets, nextLine and error.

error only return io.EOF goroutine not allowed. read files must be sync
*/
func (d *UrbanDictionaryReader) ReadNextNLines(n int) (datasets []UrbanDictionaryDatasets, nextLine int, err error) {
	d.init()
	d.mutex.Lock()
	datasets = make([]UrbanDictionaryDatasets, n)
	var invalidRow int64 = 0
	var rec []string
	for i := 0; i < n; i++ {

		rec, err = d.reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				datasets = datasets[:i]
				break
			}
			invalidRow++
			continue
		}
		datasets[i] = UrbanDictionaryDatasets{
			WordId:     rec[0],
			Word:       EscapeXSSandBackSlash(rec[1]),
			UpVotes:    rec[2],
			DownVotes:  rec[3],
			Author:     EscapeXSSandBackSlash(rec[4]),
			Definition: EscapeXSSandBackSlash(rec[5]),
		}
	}
	d.ErrorCount += invalidRow
	d.index += int64(n)
	nextLine = int(d.index) + 1
	if !errors.Is(err, io.EOF) {
		err = nil
	}
	d.mutex.Unlock()
	return
}

func LoadUrbanDictionaryDatasets() []UrbanDictionaryDatasets {
	//getting filedir fullpath
	_, fileDir, _, _ := runtime.Caller(0)
	//getting dir and concat with datasets
	f, err := os.OpenFile(fileDir[:len(fileDir)-len(`urbandict.go`)]+`urbandict-word-defs.csv`, os.O_RDONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	ud := make([]UrbanDictionaryDatasets, 0)
	rd := csv.NewReader(f)
	rd.Read()
	var rec []string
	invalidRow := 0
	for !errors.Is(err, io.EOF) {
		rec, err = rd.Read()
		if err != nil {
			invalidRow++
			continue
		}
		ud = append(ud, UrbanDictionaryDatasets{
			WordId:     rec[0],
			Word:       EscapeXSSandBackSlash(rec[1]),
			UpVotes:    rec[2],
			DownVotes:  rec[3],
			Author:     EscapeXSSandBackSlash(rec[4]),
			Definition: EscapeXSSandBackSlash(rec[5]),
		})
	}
	log.Printf("there are %v invalid row \n", invalidRow)
	return ud
}
func LoadUrbanDictionaryDatasetsDatasetsChan(ch chan<- UrbanDictionaryDatasets) {
	//getting filedir fullpath
	_, fileDir, _, _ := runtime.Caller(0)
	//getting dir and concat with datasets
	f, err := os.OpenFile(fileDir[:len(fileDir)-len(`urbandict.go`)]+`urbandict-word-defs.csv`, os.O_RDONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	rd := csv.NewReader(f)
	//skip first line which is header name of the datasets
	rd.Read()
	var rec []string
	invalidRow := 0
	for !errors.Is(err, io.EOF) {
		rec, err = rd.Read()
		if err != nil {
			invalidRow++
			continue
		}
		ch <- UrbanDictionaryDatasets{
			WordId:     rec[0],
			Word:       EscapeXSSandBackSlash(rec[1]),
			UpVotes:    rec[2],
			DownVotes:  rec[3],
			Author:     EscapeXSSandBackSlash(rec[4]),
			Definition: EscapeXSSandBackSlash(rec[5]),
		}
	}
	close(ch)
	log.Printf("there are %v invalid row \n", invalidRow)
}

func EscapeXSSandBackSlash(s string) string {
	return S.Replace(S.XSS(s), `\`, `\\`)
}
