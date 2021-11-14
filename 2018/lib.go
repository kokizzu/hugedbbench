package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const KEY_InsertUsersItems = `InsertUsersItems`
const KEY_UpdateItemsAmounts = `UpdateItemsAmounts`
const KEY_SearchRelsAddBonds = `SearchRelsAddBonds`
const KEY_RandomSearchItems = `RandomSearchItems`

var wg = sync.WaitGroup{}
var PREFIX string

const INSERT_USERS = 50 * 1000 / SLOW_DATABASE
const RELS_PER_USER = 20 / SLOW_DATABASE
const SEARCH_USER = 100 * 1000 / SLOW_DATABASE
const ITEM_PER_SEARCH = 3
const SEARCH_RELS = 200 * 1000 / SLOW_DATABASE
const SEARCH_ITEM = 100 * 1000
const PROGRESS_TICK = 20

var item_list = []string{`GOLD`, `WOOD`, `OIL`, `IRON`, `ORE`, `WATER`, `FOOD`}
var db *sql.DB
var HHKEY = [32]byte{}
var secMap = map[string]float64{}

var USERS_INS = int64(0)
var USERS_SEL = int64(0)
var ITEMS_INS = int64(0)
var ITEMS_SEL = int64(0)
var ITEMS_UPD = int64(0)
var RELS_INS = int64(0)
var RELS_SEL = int64(0)
var RELS_UPD = int64(0)
var ITEMS_LIST = int64(0)

func Report() {
	wg.Add(1)
	go func() {
		time.Sleep(3 * time.Second)
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("USERS CR    : %7d / %7d \n", USERS_INS, USERS_SEL)
	fmt.Printf("ITEMS CRU   : %7d / %7d + %7d / %d \n", ITEMS_INS, ITEMS_SEL, ITEMS_LIST, ITEMS_UPD)
	fmt.Printf("RELS  CRU   : %7d / %7d / %d \n", RELS_INS, RELS_SEL, RELS_UPD)
	fmt.Printf("SLOW FACTOR : %d \n", SLOW_DATABASE)
	ins := secMap[KEY_InsertUsersItems] / float64(USERS_INS+ITEMS_INS+RELS_INS)
	rea := secMap[KEY_RandomSearchItems] / float64(ITEMS_LIST)
	upd := (secMap[KEY_UpdateItemsAmounts]/float64(ITEMS_UPD) + secMap[KEY_SearchRelsAddBonds]/float64(RELS_UPD)) / 2
	fmt.Printf("CRU µs/rec  : %.2f / %.2f / %.2f\n", ins, rea, upd)
}

func init() {
	rand.Seed(1) // TODO: make concurrent
}

func UniqString(n int) string {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(n))
	return fmt.Sprintf("%d@%x.com", n, bs)
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func NoError(err error) bool {
	if err != nil {
		fmt.Printf("[%s] %s\n", PREFIX, err)
		return false
	}
	return true
}
func NoError2(err error, query string) bool {
	if err != nil {
		fmt.Printf("[%s] %s: %s\n", PREFIX, err, query)
		return false
	}
	return true
}

func RunThread(name string, num int, lambda func(num int)) {
	wg.Add(1)
	if name != KEY_InsertUsersItems {
		time.Sleep(DELAY_SEARCH * time.Second)
	}
	fmt.Printf("[%s] %s started..\n", PREFIX, name)
	start := time.Now()
	tick := num / PROGRESS_TICK

	for i := 1; i <= num; i++ {
		lambda(i)
		if i%tick == 0 {
			elapsed := time.Since(start)
			fmt.Printf("[%s] %s (%d, %d%%) took %.2fs (%.2f µs/op)\n", PREFIX, name, i, i*100/num, elapsed.Seconds(), float64(elapsed.Nanoseconds())/1000/float64(i))
			secMap[name] = float64(elapsed.Nanoseconds() / 1000)
		}
	}
	fmt.Printf("[%s] %s completed..\n", PREFIX, name)
	wg.Done()
}

func DdlNoError(query string) {
	fmt.Println(query)
	_, err := db.Exec(query)
	NoError(err)
}
func DdlPanicIf(query string) {
	fmt.Println(query)
	_, err := db.Exec(query)
	PanicIf(err)
}
