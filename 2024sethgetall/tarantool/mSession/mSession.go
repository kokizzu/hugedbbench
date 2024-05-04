package mSession

import (
	. "github.com/kokizzu/gotro/D/Tt"
)

// custom struct

type Session struct {
	Id         uint64
	Email      string
	Permission string
	ExpiredAt  int64
	SessionKey uint64
}

const TableSessions = `sessions`

var Tables = map[TableName]*TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TableSessions: {
		Fields: []Field{
			{IdCol, Unsigned},
			{`email`, String},
			{`permission`, String},
			{`expiredAt`, Integer},
			{`sessionKey`, String},
		},
		AutoIncrementId: false,
		Engine:          Memtx,
		Unique1:         `sessionKey`,
	},
}
