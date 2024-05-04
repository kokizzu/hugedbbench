package rqPoints

import (
	"github.com/kokizzu/gotro/L"
)

// FindNearestPoints tarantool can't search by distance
func (p *PointsSg) FindNearestPoints(_ float64, limit int64) []PointsSg {
	var rows []PointsSg
	// p.Coord is lat, long
	res, err := p.Adapter.Select(p.SpaceName(), p.SpatialIndexCoord(), 0, uint32(limit), 11, p.Coord)
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
		return rows
	}
	for _, row := range res {
		item := PointsSg{}
		if row, ok := row.([]any); ok {
			rows = append(rows, *item.FromArray(row))
		}
	}
	return rows
}
