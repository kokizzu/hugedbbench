package mTest

import (
	. "github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool"
)

func Migrate(taran *Adapter) {
	taran.MigrateTables(Tables)
}

func ConnectTarantool() *tarantool.Connection {
	taran, err := tarantool.Connect(`localhost:3301`, tarantool.Opts{})
	L.PanicIf(err, `failed to connect to tarantool`)
	return taran
}
