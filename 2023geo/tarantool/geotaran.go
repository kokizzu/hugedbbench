package main

import (
	"errors"
	"math"

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
		// decorate with distance
		deco := make([]any, 0, len(res))
		for _, r := range res {
			lat2 := r.Coord[0].(float64)
			long2 := r.Coord[1].(float64)
			w := lat - lat2
			h := long - long2
			d := math.Sqrt(w*w+h*h) * geo.DegToMeter
			deco = append(deco, []any{r.Id, lat2, long2, d})
		}
		_ = deco
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
