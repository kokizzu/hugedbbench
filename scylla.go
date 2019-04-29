package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"math/rand"
	"sync/atomic"
	"time"
)

const SLOW_DATABASE = 20 // set to 1 for fast database
const DELAY_SEARCH = 2   // set to seconds required for first insert passed

var highest_user_id = 0
var sess *gocql.Session

func main() {
	var err error
	cl := gocql.NewCluster("127.0.0.1")
	cl.Keyspace = `b1`
	cl.Consistency = gocql.One
	cl.Timeout = 5 * time.Second
	sess, err = cl.CreateSession()
	PanicIf(err)
	PREFIX = `Scylla`
	InitTables()
	go RunThread(KEY_InsertUsersItems, INSERT_USERS, InsertUsersItems)
	go RunThread(KEY_UpdateItemsAmounts, SEARCH_USER, SearchUserUpdateItemsAmounts)
	go RunThread(KEY_SearchRelsAddBonds, SEARCH_RELS, SearchRelsAddBonds)
	go RunThread(KEY_RandomSearchItems, SEARCH_ITEM, RandomSearchItems)

	Report()
	defer sess.Close()
}

func RandomSearchItems(n int) {
	query := `SELECT typ, amount FROM items WHERE user_id = ?`
	user_id := rand.Int() % highest_user_id
	iter := sess.Query(query, user_id).Iter()
	var str string
	var am int
	for iter.Scan(&str, &am) {
		atomic.AddInt64(&ITEMS_LIST, 1)
	}
	if err := iter.Close(); NoError(err) {
	}
}

func SearchRelsAddBonds(n int) {
	query := `SELECT user_lo, bond FROM rels WHERE user_hi = ? ALLOW FILTERING`
	if RELS_SEL%2 == 0 {
		query = `SELECT user_hi, bond FROM rels WHERE user_lo = ? ALLOW FILTERING`
	}
	user_id := rand.Int() % highest_user_id
	iter := sess.Query(query, user_id).Iter()
	var friend_id, bond int
	for iter.Scan(&friend_id, &bond) {
		atomic.AddInt64(&RELS_SEL, 1)
		//if RELS_SEL%2 == 0 {
		user_hi := user_id
		user_lo := friend_id
		if user_hi < user_lo {
			user_hi, user_lo = user_lo, user_hi
		}
		bond += 1
		query = `UPDATE rels SET bond = ? WHERE user_lo = ? AND user_hi = ?`
		err := sess.Query(query, bond, user_lo, user_hi).Exec()
		if NoError2(err, query) {
			atomic.AddInt64(&RELS_UPD, 1)
		}
		//}
	}
	if err := iter.Close(); NoError2(err, query) {
	}
}

func SearchUserUpdateItemsAmounts(n int) {
	query := `SELECT id FROM users WHERE uniq = ?`
	r := 1 + rand.Int()%highest_user_id
	uniq := UniqString(r)
	user_id := 0
	err := sess.Query(query, uniq).Scan(&user_id)
	if NoError2(err, query) {
		atomic.AddInt64(&USERS_SEL, 1)
	}

	// random 3 item
	idxs := [ITEM_PER_SEARCH]int{}
	dummy := [ITEM_PER_SEARCH]int{}
	for t := 0; t < ITEM_PER_SEARCH; {
		idx := rand.Int() % len(item_list)
		x := 0
		for x < t {
			if idxs[x] == idx {
				idx = rand.Int() % len(item_list)
				continue
			}
			x++
		}
		idxs[t] = idx
		t++
	}
	query = `SELECT user_id, amount FROM items WHERE user_id = ? AND typ = ?`
	var amount int // dummy
	for k, idx := range idxs {
		err = sess.Query(query, user_id, item_list[idx]).Scan(&dummy[k], &amount)
		if NoError2(err, query) {
			atomic.AddInt64(&ITEMS_SEL, 1)
			dummy[k] = amount
		}
	}

	// without transaction to make it fair with other databases
	for k, idx := range idxs {
		query = fmt.Sprintf(`UPDATE items SET amount = %d WHERE user_id = %d AND typ = '%s'`, dummy[k]+rand.Int()%10-3, user_id, item_list[idx])
		err := sess.Query(query).Exec()
		if NoError2(err, query) {
			atomic.AddInt64(&ITEMS_UPD, 1)
		}
	}
}

func InsertUsersItems(user_id int) {
	query := `INSERT INTO users(id, uniq, created_at) VALUES (?, ?, ?)`
	uniq := UniqString(user_id)
	err := sess.Query(query, user_id, uniq, time.Now().Unix()).Exec()
	if NoError(err) {
		atomic.AddInt64(&USERS_INS, 1)
	}
	if highest_user_id < user_id { // TODO: should be atomic
		highest_user_id = user_id
	}
	for _, item := range item_list {
		query = `INSERT INTO items (user_id, typ, amount, created_at) VALUES(?, ?, ?, ?)`
		err := sess.Query(query, user_id, item, rand.Int()%100, time.Now().Unix()).Exec()
		if NoError(err) {
			atomic.AddInt64(&ITEMS_INS, 1)
		}
	}
	delta := INSERT_USERS / PROGRESS_TICK
	if user_id > delta {
		for r := 0; r < RELS_PER_USER; r++ {
			lo := user_id - delta + rand.Int()%delta
			hi := user_id
			query = `INSERT INTO rels(user_lo,user_hi, bond,created_at) VALUES(?,?, 0,?)`
			err = sess.Query(query, lo, hi, time.Now().Unix()).Exec()
			if NoError(err) {
				atomic.AddInt64(&RELS_INS, 1)
			}
		}
	}
}

func InitTables() {
	ScyNoError(`DROP TABLE IF EXISTS users`)
	ScyNoError(`DROP TABLE IF EXISTS items`)
	ScyNoError(`DROP TABLE IF EXISTS rels`)
	ScyPanicIf(`CREATE TABLE users(
   id INT
	, uniq TEXT
	, created_at TIMESTAMP 
	, PRIMARY KEY(uniq)
)`)
	//ScyPanicIf(`CREATE UNIQUE INDEX IF NOT EXISTS users__uniq ON users(uniq)`)
	ScyPanicIf(`CREATE TABLE items(
    user_id INT -- REFERENCES users(id)
	, typ TEXT
	, amount INT 
	, created_at TIMESTAMP 
	, PRIMARY KEY(user_id, typ)
)`)
	//ScyPanicIf(`CREATE UNIQUE INDEX IF NOT EXISTS items__user_id__typ ON items(user_id, typ)`)
	ScyPanicIf(`CREATE TABLE rels(
   user_lo INT -- REFERENCES users(id)
	, user_hi INT -- REFERENCES users(id)
	, bond INT 
	, created_at TIMESTAMP 
	, PRIMARY KEY(user_lo, user_hi)
)`)
	//ScyPanicIf(`CREATE INDEX IF NOT EXISTS rels__user_lo ON rels(user_lo)`)
	//ScyPanicIf(`CREATE INDEX IF NOT EXISTS rels__user_hi ON rels(user_hi)`)
}

func ScyPanicIf(query string) {
	fmt.Println(query)
	err := sess.Query(query).Exec()
	PanicIf(err)
}

func ScyNoError(query string) {
	fmt.Println(query)
	err := sess.Query(query).Exec()
	NoError(err)
}
