package rqTest

import (
	"github.com/kokizzu/gotro/A"
)

func (t *TestTable2) FindWhereInStruct(ids []uint64) (res []TestTable2) {
	query := `
SELECT ` + t.SqlId() + `
, ` + t.SqlContent() + `
FROM ` + t.SqlTableName() + `
WHERE ` + t.SqlId() + ` IN (` + A.UIntJoin(ids, `,`) + `)
`
	t.Adapter.QuerySql(query, func(row []any) {
		row2 := TestTable2{}
		row2.FromArray(row)
		res = append(res, row2)
	})
	return
}

func (t *TestTable2) FindWhereInArray(ids []uint64) (res [][]any) {
	query := `
SELECT ` + t.SqlId() + `
, ` + t.SqlContent() + `
FROM ` + t.SqlTableName() + `
WHERE ` + t.SqlId() + ` IN (` + A.UIntJoin(ids, `,`) + `)
`
	t.Adapter.QuerySql(query, func(row []any) {
		res = append(res, row)
	})
	return
}
