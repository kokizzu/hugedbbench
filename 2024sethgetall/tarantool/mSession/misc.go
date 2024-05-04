package mSession

import (
	. "github.com/kokizzu/gotro/D/Tt"
	"github.com/tarantool/go-tarantool/v2"
)

func Migrate(taran *Adapter) {
	taran.MigrateTables(Tables)
}

func ConnectTarantool() *tarantool.Connection {
	return Connect1(`127.0.0.1`, `3301`, `user1`, `password1`)
}
