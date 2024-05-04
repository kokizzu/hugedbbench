package main

import (
	"context"
	"time"

	"hugedbbench/2024sethgetall/testcase"

	"github.com/kokizzu/gotro/B"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/rueian/rueidis"
)

// garnet: dotnet restore && cd main/GarnetServer && dotnet run -c Release -f net8.0 # use port 3278
// dragonflydb: docker run --network=host --ulimit memlock=-1 docker.dragonflydb.io/dragonflydb/dragonfly # must disable clientCaching
// keydb: docker run --name some-keydb2 -p 6379:6379 -d eqalpha/keydb keydb-server /etc/keydb/keydb.conf --server-threads 4
// kvrocks: docker run -it -p 6379:6666 apache/kvrocks

func main() {
	cli, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{`127.0.0.1:6379`},
		Password:    `kl234j23095125125125`,
		//AlwaysRESP2:  true, // for garnet
		//DisableCache: true, // for dragonflydb, kvrocks, garnet
	})
	L.PanicIf(err, `rueidis.NewClient`)
	defer cli.Close()

	ctx := context.Background()

	// SET
	keys := testcase.RunInsert(`SET`, func(i int) string {
		sessionKey, byt := testcase.CreateSessionByt(i)
		b := cli.B().Set().Key(sessionKey).Value(string(byt)).ExSeconds(testcase.ExpireSec).Build()
		resp := cli.Do(ctx, b)
		L.IsError(resp.Error(), `failed to SET`, sessionKey)
		return sessionKey
	})

	// GET
	testcase.RunSearch(`GET`, keys, func(sessionKey string) bool {
		b := cli.B().Get().Key(sessionKey).Build()
		byt, err := cli.Do(ctx, b).AsBytes()
		if L.IsError(err, `failed to GET`, sessionKey) {
			return false
		}
		session, valid := testcase.ReadSessionByt(sessionKey, byt)
		return session.Id != 0 && valid
	})
	// DEL
	start := time.Now()
	for _, sessionKey := range keys {
		b := cli.B().Del().Key(sessionKey).Build()
		resp := cli.Do(ctx, b)
		L.IsError(resp.Error(), `failed to DEL`, sessionKey)
	}
	L.TimeTrack(start, `DEL 10k user session`)

	// HSET
	keys = testcase.RunInsert(`HSET+TTL`, func(i int) string {
		sessionKey, session := testcase.CreateSession(i)
		b2 := cli.B().Hset().Key(sessionKey).FieldValue().
			FieldValue(`id`, I.ToS(session.Id)).
			FieldValue(`email`, session.Email)
		for k, ok := range session.Permission {
			if ok {
				b2 = b2.FieldValue(k, B.ToS(ok))
			}
		}
		b := b2.Build()
		resp := cli.Do(ctx, b)
		L.IsError(resp.Error(), `failed to HSET`, sessionKey)
		b = cli.B().Expire().Key(sessionKey).Seconds(testcase.ExpireSec).Build()
		resp = cli.Do(ctx, b)
		L.IsError(resp.Error(), `failed to EXPIRE`, sessionKey)
		return sessionKey
	})

	// HGETALL
	testcase.RunSearch(`HGETALL`, keys, func(sessionKey string) bool {
		b := cli.B().Hgetall().Key(sessionKey).Build()
		rows, err := cli.Do(ctx, b).AsStrMap()
		if L.IsError(err, `failed to HGETALL`, sessionKey) {
			return false
		}
		session := testcase.Session{Permission: map[string]bool{}}
		for key, value := range rows {
			switch key {
			case `id`:
				session.Id = S.ToI(value)
			case `email`:
				session.Email = value
			default:
				if testcase.PossiblePerm[key] {
					session.Permission[key] = value == `true`
				}
			}
		}
		return session.Id != 0
	})

	// DEL
	start = time.Now()
	for _, sessionKey := range keys {
		b := cli.B().Del().Key(sessionKey).Build()
		resp := cli.Do(ctx, b)
		L.IsError(resp.Error(), `failed to DEL`, sessionKey)
	}
	L.TimeTrack(start, `DEL 10k user session`)

}
