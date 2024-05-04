package main

import (
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

	keys := testcase.RunInsert(`INSERT`, func(i int) string {
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
		return sessionKey
	})

	testcase.RunSearch(`SELECT`, keys, func(sessionKey string) bool {
		sess := rqSession.NewSessions(tt)
		sess.SessionKey = sessionKey

		if !sess.FindBySessionKey() {
			return false
		}
		session := sess.ToSession()
		if sess.ExpiredAt > fastime.Now().Unix() {
			// expired
			session.Id = 0
		}
		_ = session
		return true
	})
}
