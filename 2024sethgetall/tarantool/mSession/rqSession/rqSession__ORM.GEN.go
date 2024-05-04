package rqSession

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2024sethgetall/tarantool/mSession`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type rqSession__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type rqSession__ORM.GEN.go
// Sessions DAO reader/query struct
type Sessions struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
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

// UniqueIndexSessionKey return unique index name
func (s *Sessions) UniqueIndexSessionKey() string { //nolint:dupl false positive
	return `sessionKey`
}

// FindBySessionKey Find one by SessionKey
func (s *Sessions) FindBySessionKey() bool { //nolint:dupl false positive
	res, err := s.Adapter.Connection.Do(
		tarantool.NewSelectRequest(s.SpaceName()).
		Index(s.UniqueIndexSessionKey()).
		Offset(0).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(A.X{s.SessionKey}),
	).Get()
	if L.IsError(err, `Sessions.FindBySessionKey failed: `+s.SpaceName()) {
		return false
	}
	if len(res) == 1 {
		if row, ok := res[0].([]any); ok {
			s.FromArray(row)
			return true
		}
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
func (s *Sessions) ToUpdateArray() *tarantool.Operations { //nolint:dupl false positive
	return tarantool.NewOperations().
		Assign(0, s.Id).
		Assign(1, s.Email).
		Assign(2, s.Permission).
		Assign(3, s.ExpiredAt).
		Assign(4, s.SessionKey)
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
	return A.X{
		s.Id,         // 0
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
	res, err := s.Adapter.Connection.Do(
		tarantool.NewSelectRequest(s.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).Get()
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows
	}
	for _, row := range res {
		item := Sessions{}
		row, ok := row.([]any)
		if ok {
			rows = append(rows, *item.FromArray(row))
		}
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (s *Sessions) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	resp, err := s.Adapter.Connection.Do(
		tarantool.NewSelectRequest(s.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).GetResponse()
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	res, err := resp.Decode()
	if L.IsError(err, `Sessions.FindOffsetLimit failed: `+s.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	rows = make([]A.X, len(res))
	for _, row := range res {
		row, ok := row.([]any)
		if ok {
			rows = append(rows, row)
		}
	}
	return rows, Tt.QueryMetaFrom(resp, nil)
}

// Total count number of rows
func (s *Sessions) Total() int64 { //nolint:dupl false positive
	rows := s.Adapter.CallBoxSpace(s.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// SessionsFieldTypeMap returns key value of field name and key
var SessionsFieldTypeMap = map[string]Tt.DataType { //nolint:dupl false positive
	`id`:         Tt.Unsigned,
	`email`:      Tt.String,
	`permission`: Tt.String,
	`expiredAt`:  Tt.Integer,
	`sessionKey`: Tt.String,
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

