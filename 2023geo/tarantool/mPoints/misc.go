package mPoints

import (
	"fmt"

	. "github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool"
)

func Migrate(taran *Adapter) {
	taran.MigrateTables(Tables)
}

func ConnectTarantool() *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%d`,
		`127.0.0.1`,
		3301,
	)
	taran, err := tarantool.Connect(hostPort, tarantool.Opts{
		User: `user`,
		Pass: `password`,
	})
	L.PanicIf(err, `failed to connect to tarantool`)
	return taran
}
