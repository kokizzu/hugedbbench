package rqSession

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"hugedbbench/2024sethgetall/tarantool/mSession"

	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

// Sessions DAO reader/query struct
//
//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type rqSession__ORM.GEN.go
type Sessions struct {
	Adapter    *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id         uint64
	Email      string
	Permission string
	ExpiredAt  int64
	SessionKey string
}

// NewSessions create new ORM reader/query object
func NewSessions(adapter *Tt.Adapter) *Sessions {
	return &Sessions{Adapter: adapter}
}

// SpaceName returns full package and table name
func (s *Sessions) SpaceName() string { //nolint:dupl false positive
	return string(mSession.TableSessions) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (s *Sessions) SqlTableName() string { //nolint:dupl false positive
	return `"sessions"`
}

func (s *Sessions) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (s *Sessions) FindById() bool { //nolint:dupl false positive
	res, err := s.Adapter.Select(s.SpaceName(), s.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{s.Id})
	if L.IsError(err, `Sessions.FindById failed: `+s.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		s.FromArray(rows[0])
		return true
	}
	return false
}

// UniqueIndexSessionKey return unique index name
func (s *Sessions) UniqueIndexSessionKey() string { //nolint:dupl false positive
	return `sessionKey`
}

// FindBySessionKey Find one by SessionKey
func (s *Sessions) FindBySessionKey() bool { //nolint:dupl false positive
	res, err := s.Adapter.Select(s.SpaceName(), s.UniqueIndexSessionKey(), 0, 1, tarantool.IterEq, A.X{s.SessionKey})
	if L.IsError(err, `Sessions.FindBySessionKey failed: `+s.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		s.FromArray(rows[0])
		return true
	}
	return false
}

// SqlSelectAllFields generate Sql select fields
func (s *Sessions) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "email"
	, "permission"
	, "expiredAt"
	, "sessionKey"
	`
}

// SqlSelectAllUncensoredFields generate Sql select fields
func (s *Sessions) SqlSelectAllUncensoredFields() string { //nolint:dupl false positive
	return ` "id"
	, "email"
	, "permission"
	, "expiredAt"
	, "sessionKey"
	`
}

// ToUpdateArray generate slice of update command
func (s *Sessions) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, s.Id},
		A.X{`=`, 1, s.Email},
		A.X{`=`, 2, s.Permission},
		A.X{`=`, 3, s.ExpiredAt},
		A.X{`=`, 4, s.SessionKey},
	}
}

// IdxId return name of the index
func (s *Sessions) IdxId() int { //nolint:dupl false positive
	return 0
}

// SqlId return name of the column being indexed
func (s *Sessions) SqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxEmail return name of the index
func (s *Sessions) IdxEmail() int { //nolint:dupl false positive
	return 1
}

// SqlEmail return name of the column being indexed
func (s *Sessions) SqlEmail() string { //nolint:dupl false positive
	return `"email"`
}

// IdxPermission return name of the index
func (s *Sessions) IdxPermission() int { //nolint:dupl false positive
	return 2
}

// SqlPermission return name of the column being indexed
func (s *Sessions) SqlPermission() string { //nolint:dupl false positive
	return `"permission"`
}

// IdxExpiredAt return name of the index
func (s *Sessions) IdxExpiredAt() int { //nolint:dupl false positive
	return 3
}

// SqlExpiredAt return name of the column being indexed
func (s *Sessions) SqlExpiredAt() string { //nolint:dupl false positive
	return `"expiredAt"`
}

// IdxSessionKey return name of the index
func (s *Sessions) IdxSessionKey() int { //nolint:dupl false positive
	return 4
}

// SqlSessionKey return name of the column being indexed
func (s *Sessions) SqlSessionKey() string { //nolint:dupl false positive
	return `"sessionKey"`
}

// ToArray receiver fields to slice
func (s *Sessions) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if s.Id != 0 {
		id = s.Id
	}
	return A.X{
		id,
		s.Email,      // 1
		s.Permission, // 2
		s.ExpiredAt,  // 3
		s.SessionKey, // 4
	}
}

// FromArray convert slice to receiver fields
func (s *Sessions) FromArray(a A.X) *Sessions { //nolint:dupl false positive
	s.Id = X.ToU(a[0])
	s.Email = X.ToS(a[1])
	s.Permission = X.ToS(a[2])
	s.ExpiredAt = X.ToI(a[3])
	s.SessionKey = X.ToS(a[4])
	return s
}

// FromUncensoredArray convert slice to receiver fields
func (s *Sessions) FromUncensoredArray(a A.X) *Sessions { //nolint:dupl false positive
	s.Id = X.ToU(a[0])
	s.Email = X.ToS(a[1])
	s.Permission = X.ToS(a[2])
	s.ExpiredAt = X.ToI(a[3])
	s.SessionKey = X.ToS(a[4])
	return s
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (s *Sessions) FindOffsetLimit(offset, limit uint32, idx string) []Sessions { //nolint:dupl false positive
	var rows []Sessions
	res, err := s.Adapter.Select(s.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows
	}
	for _, row := range res.Tuples() {
		item := Sessions{}
		rows = append(rows, *item.FromArray(row))
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (s *Sessions) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	res, err := s.Adapter.Select(s.SpaceName(), idx, offset, limit, tarantool.IterAll, A.X{})
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows, Tt.QueryMetaFrom(res, err)
	}
	tuples := res.Tuples()
	rows = make([]A.X, len(tuples))
	for z, row := range tuples {
		rows[z] = row
	}
	return rows, Tt.QueryMetaFrom(res, nil)
}

// Total count number of rows
func (s *Sessions) Total() int64 { //nolint:dupl false positive
	rows := s.Adapter.CallBoxSpace(s.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// SessionsFieldTypeMap returns key value of field name and key
var SessionsFieldTypeMap = map[string]Tt.DataType{ //nolint:dupl false positive
	`id`:         Tt.Unsigned,
	`email`:      Tt.String,
	`permission`: Tt.String,
	`expiredAt`:  Tt.Integer,
	`sessionKey`: Tt.String,
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
