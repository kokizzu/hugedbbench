package mPoints

import (
	"github.com/kokizzu/gotro/A"
	. "github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/X"
)

// custom struct

type PointsSg struct {
	Id    uint64
	Coord []float64
}

func (u *PointsSg) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if u.Id != 0 {
		id = u.Id
	}
	return A.X{
		id,
		u.Coord, // 1
	}
}

func (u *PointsSg) FromArray(a A.X) *PointsSg { //nolint:dupl false positive
	u.Id = X.ToU(a[0])
	u.Coord = X.ToAF(a[1])
	return u
}

func (u *PointsSg) ToMapFromSlice(row []any) map[string]any {
	return map[string]any{
		IdCol:  row[0],
		`lat`:  row[1],
		`long`: row[2],
	}
}

const TablePointsSg = `points_sg`

var Tables = map[TableName]*TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TablePointsSg: {
		Fields: []Field{
			{IdCol, Unsigned},
			{`coord`, Array},
		},
		AutoIncrementId: true,
		Engine:          Vinyl,
	},
	// TODO: add support for rtree
	// https://www.tarantool.io/en/doc/latest/concepts/data_model/indexes/
	// to be fair, not making an index on "content"
}
