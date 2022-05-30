package datasets

import (
	"log"
	"testing"
)

func TestLoadUrbandict(t *testing.T) {
	// dir, _ := os.Getwd()
	ms := LoadUrbanDictionaryDatasets()
	if ms == nil {
		t.Error(`empty datasets`)
	}
	log.Println(len(ms))
}

func TestLoadUrbanReader(t *testing.T) {
	r := UrbanDictionaryReader{SkipHeader: true}
	ud, _, err := r.ReadNextNLines(50)
	if err != nil {
		panic(err)
	}
	log.Println(ud)
}
