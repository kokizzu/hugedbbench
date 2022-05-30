package datasets

import (
	"log"
	"testing"
)

func TestLoadMassive(t *testing.T) {
	// dir, _ := os.Getwd()
	ms := ReadMassiveDatasets()
	if ms == nil {
		t.Error(`empty datasets`)
	}
	log.Println(len(ms))
}
