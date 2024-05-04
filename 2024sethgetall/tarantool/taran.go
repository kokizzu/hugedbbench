package main

import (
	"fmt"
	"math/rand"
	"time"

	"hugedbbench/2023geo/tarantool/mPoints"
	"hugedbbench/2024sethgetall/tarantool/mSession"
	"hugedbbench/2024sethgetall/tarantool/mSession/rqSession"
	"hugedbbench/2024sethgetall/tarantool/mSession/wcSession"
	"hugedbbench/2024sethgetall/testcase"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kpango/fastime"
	"golang.org/x/exp/maps"
)

func main() {

	tt := &Tt.Adapter{Connection: mPoints.ConnectTarantool(), Reconnect: mPoints.ConnectTarantool}
	_, err := tt.Ping()
	L.PanicIf(err, `tt.Ping`)

	// clear databae
	mSession.Migrate(tt)

	tt.TruncateTable(mSession.TableSessions)

	keys := make([]string, 0, testcase.UserCount)

	start := time.Now()
	for i := range testcase.UserCount {
		sessionKey, session := testcase.CreateSession(i)
		sess := wcSession.NewSessionsMutator(tt)
		sess.Id = uint64(i)
		sess.Email = session.Email
		sess.Permission = A.StrJoin(maps.Keys(session.Permission), ` `)
		sess.SessionKey = sessionKey
		sess.ExpiredAt = fastime.Now().Unix() + testcase.ExpireSec
		if !sess.DoInsert() {
			L.Print(`failed to insert`, sess)
		}
		keys = append(keys, sessionKey)
	}
	ms := L.TimeTrack(start, `INSERT 10k user session`)
	fmt.Printf("%.0f rps\n", testcase.UserCount/ms*1000)

	start = time.Now()
	failCount := 0
	for range testcase.RequestCount {
		i := rand.Intn(len(keys)) // assume this is per request
		sessionKey := keys[i]

		sess := rqSession.NewSessions(tt)
		sess.SessionKey = sessionKey

		if !sess.FindBySessionKey() {
			failCount++
			continue
		}
		_ = sess.ToSession()
	}
	ms = L.TimeTrack(start, `SELECT 10k 20x user session`)
	fmt.Printf("%.0f rps, failed: %d\n", testcase.RequestCount/ms*1000, failCount)
}
