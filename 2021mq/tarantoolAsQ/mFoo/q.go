package mFoo

import "github.com/kokizzu/gotro/D/Tt"

const (
	TableFoo Tt.TableName = `foo`
	Id                    = `id`
	When                  = `when`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TableFoo: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{When, Tt.Unsigned},
		},
		AutoIncrementId: true,
		Engine:          Tt.Memtx,
		//Engine:        Tt.Vinyl,
	},
}
