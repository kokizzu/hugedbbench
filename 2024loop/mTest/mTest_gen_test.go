package mTest

import (
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
)

// only need to do once before compile
func TestGenerateOrm(t *testing.T) {
	Tt.GenerateOrm(Tables, false)
	t.SkipNow()
}
func TestMigrate(t *testing.T) {
	taran := &Tt.Adapter{Connection: ConnectTarantool(), Reconnect: ConnectTarantool}
	_, err := taran.Ping()
	L.PanicIf(err, `taran.Ping`)
	Migrate(taran)
}
