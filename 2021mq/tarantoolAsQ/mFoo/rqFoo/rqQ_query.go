package rqFoo

import "github.com/kokizzu/gotro/I"

func (s *Foo) FindGreaterThan(id, limit int64) (res []*Foo) {
	query := `
SELECT ` + s.sqlSelectAllFields() + `
FROM ` + s.sqlTableName() + `
WHERE ` + s.sqlId() + `  > ` + I.ToS(id) +`
ORDER BY ` + s.sqlId() + `
LIMIT `+I.ToS(limit) // note: for string, use S.Z or S.XSS to prevent SQL injection
	s.Adapter.QuerySql(query, func(row []interface{}) {
		obj := &Foo{}
		obj.FromArray(row)
		res = append(res, obj)
	})
	return
}
