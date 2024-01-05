package mTest

import (
	"github.com/kokizzu/gotro/A"
	. "github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/X"
)

// custom struct

type TestTable2 struct {
	Id      uint64
	Content string
}

func (u *TestTable2) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if u.Id != 0 {
		id = u.Id
	}
	return A.X{
		id,
		u.Content, // 1
	}
}

func (u *TestTable2) FromArray(a A.X) *TestTable2 { //nolint:dupl false positive
	u.Id = X.ToU(a[0])
	u.Content = X.ToS(a[1])
	return u
}

func (u *TestTable2) ToMapFromSlice(row []any) map[string]any {
	return map[string]any{
		IdCol:     row[0],
		`content`: row[1],
	}
}

const TableTestTable2 = `test_table2`

var Tables = map[TableName]*TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TableTestTable2: {
		Fields: []Field{
			{IdCol, Unsigned},
			{`content`, String},
		},
		Unique1:         `content`,
		AutoIncrementId: true,
		//Engine:          Vinyl,
		Engine: Memtx,
	},
	// to be fair, not making an index on "content"
}
