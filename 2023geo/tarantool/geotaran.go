package main

import (
	"errors"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"

	geo "hugedbbench/2023geo"
	"hugedbbench/2023geo/tarantool/mPoints"
	"hugedbbench/2023geo/tarantool/mPoints/rqPoints"
	"hugedbbench/2023geo/tarantool/mPoints/wcPoints"
)

func main() {

	tt := &Tt.Adapter{Connection: mPoints.ConnectTarantool(), Reconnect: mPoints.ConnectTarantool}
	_, err := tt.Ping()
	L.PanicIf(err, `tt.Ping`)

	const key = `SG` // per city

	// clear databae
	mPoints.Migrate(tt)

	tt.TruncateTable(mPoints.TablePointsSg)

	geo.Insert100kPoints(func(lat, long float64, id uint64) error {
		p := wcPoints.NewPointsSgMutator(tt)
		p.Id = id
		p.Coord = []any{lat, long}
		if p.DoInsert() {
			return nil
		}
		return errors.New(`DoInsert`)
	})

	geo.SearchRadius200k(func(lat, long, boxMeter float64, maxResult int64) (uint64, error) {
		p := rqPoints.NewPointsSg(tt)
		p.Coord = []any{lat, long}
		res := p.FindNearestPoints(boxMeter, maxResult)
		l := len(res)

		return uint64(l), nil
	})

	geo.MovingPoint(func(lat, long float64, id uint64) error {
		p := wcPoints.NewPointsSgMutator(tt)
		p.Id = id
		if p.FindById() {
			p.Coord = []any{lat, long}
			if !p.DoUpdateById() {
				return errors.New(`DoUpdateById`)
			}
		}
		return nil
	})
}
