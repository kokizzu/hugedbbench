package testcase

import (
	"fmt"
	"strings"
	"time"

	"github.com/kokizzu/gotro/S"
	"github.com/vmihailenco/msgpack/v5"
)

type Session struct {
	Id         int64
	Email      string
	Permission map[string]bool
}

const UserCount = 10_000
const ExpireSec = 24 * 30 * 60 * 60
const RequestCount = UserCount * 20

var PossiblePerm = map[string]bool{`user`: true, `admin`: true, `staff`: true}

func CreateSession(i int) (sessionKey string, session Session) {
	perm := map[string]bool{`user`: true}
	if i%1000 == 0 {
		perm[`admin`] = true
	}
	if i%400 == 0 {
		perm[`staff`] = true
	}
	session = Session{
		Id:         int64(i),
		Email:      fmt.Sprintf(`email%d@gmail.com`, i),
		Permission: perm,
	}
	sessionKey = `session|` + S.EncodeCB63(int64(i), 1) + `|` + S.EncodeCB63(time.Now().UnixNano(), 1)
	sessionKey += `|` + S.HashPassword(session.Email)
	return sessionKey, session
}

func CreateSessionByt(i int) (sessionKey string, byt []byte) {
	sessionKey, session := CreateSession(i)
	byt, _ = msgpack.Marshal(session)
	return sessionKey, byt
}

func ReadSession(sessionKey string, session Session) bool {
	v := strings.Split(sessionKey, `|`)
	if len(v) != 4 {
		return false
	}
	if S.HashPassword(session.Email) != v[3] {
		return false
	}
	id, ok := S.DecodeCB63[int64](v[1])
	if !ok || id != session.Id {
		return false
	}
	return true
}

func ReadSessionByt(sessionKey string, value []byte) (session Session, valid bool) {
	err := msgpack.Unmarshal(value, &session)
	if err != nil {
		return Session{}, false
	}
	if ReadSession(sessionKey, session) {
		return session, true
	}
	return Session{}, false
}
