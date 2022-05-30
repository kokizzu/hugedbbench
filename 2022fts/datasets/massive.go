package datasets

import (
	"bufio"
	"encoding/json"
	"os"
	"runtime"
)

//{"id": "0", "locale": "en-US", "partition": "test", "scenario": "alarm", "intent": "alarm_set", "utt": "wake me up at five am this week", "annot_utt": "wake me up at [time : five am] [date : this week]", "worker_id": "1"}

//go:generate gomodifytags -file massive.go -struct MassiveDatasets -add-tags json -w
type MassiveDatasets struct {
	Id        string `json:"id"`
	Locale    string `json:"locale"`
	Partition string `json:"partition"`
	Scenario  string `json:"scenario"`
	Intent    string `json:"intent"`
	Utt       string `json:"utt"`
	AnnotUtt  string `json:"annot_utt"`
	WorkerId  string `json:"worker_id"`
}

func ReadMassiveDatasets() []MassiveDatasets {

	_, fs, _, _ := runtime.Caller(0)

	f, err := os.OpenFile(fs[:len(fs)-len(`massive.go`)]+`massive_en-US.jsonl`, os.O_RDONLY, os.ModeType)
	if err != nil {
		panic(err)
	}
	bf := bufio.NewScanner(f)
	massives := make([]MassiveDatasets, 0)

	for bf.Scan() {

		massive := MassiveDatasets{}
		err := json.Unmarshal(bf.Bytes(), &massive)
		if err != nil {
			panic(err)
		}
		massives = append(massives, massive)
	}

	return massives
}
