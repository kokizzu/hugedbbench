package rqFoo

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2021mq/tarantoolAsQ/mFoo`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqFoo__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type rqFoo__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type rqFoo__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type rqFoo__ORM.GEN.go
// Foo DAO reader/query struct
type Foo struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id   uint64
	When uint64
}

// NewFoo create new ORM reader/query object
func NewFoo(adapter *Tt.Adapter) *Foo {
	return &Foo{Adapter: adapter}
}

// SpaceName returns full package and table name
func (f *Foo) SpaceName() string { //nolint:dupl false positive
	return string(mFoo.TableFoo) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (f *Foo) SqlTableName() string { //nolint:dupl false positive
	return `"foo"`
}

func (f *Foo) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (f *Foo) FindById() bool { //nolint:dupl false positive
	res, err := f.Adapter.Connection.Do(
		tarantool.NewSelectRequest(f.SpaceName()).
		Index(f.UniqueIndexId()).
		Offset(0).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(A.X{f.Id}),
	).Get()
	if L.IsError(err, `Foo.FindById failed: `+f.SpaceName()) {
		return false
	}
	if len(res) == 1 {
		if row, ok := res[0].([]any); ok {
			f.FromArray(row)
			return true
		}
	}
	return false
}

// SqlSelectAllFields generate Sql select fields
func (f *Foo) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "when"
	`
}

// SqlSelectAllUncensoredFields generate Sql select fields
func (f *Foo) SqlSelectAllUncensoredFields() string { //nolint:dupl false positive
	return ` "id"
	, "when"
	`
}

// ToUpdateArray generate slice of update command
func (f *Foo) ToUpdateArray() *tarantool.Operations { //nolint:dupl false positive
	return tarantool.NewOperations().
		Assign(0, f.Id).
		Assign(1, f.When)
}

// IdxId return name of the index
func (f *Foo) IdxId() int { //nolint:dupl false positive
	return 0
}

// SqlId return name of the column being indexed
func (f *Foo) SqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxWhen return name of the index
func (f *Foo) IdxWhen() int { //nolint:dupl false positive
	return 1
}

// SqlWhen return name of the column being indexed
func (f *Foo) SqlWhen() string { //nolint:dupl false positive
	return `"when"`
}

// ToArray receiver fields to slice
func (f *Foo) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if f.Id != 0 {
		id = f.Id
	}
	return A.X{
		id,
		f.When, // 1
	}
}

// FromArray convert slice to receiver fields
func (f *Foo) FromArray(a A.X) *Foo { //nolint:dupl false positive
	f.Id = X.ToU(a[0])
	f.When = X.ToU(a[1])
	return f
}

// FromUncensoredArray convert slice to receiver fields
func (f *Foo) FromUncensoredArray(a A.X) *Foo { //nolint:dupl false positive
	f.Id = X.ToU(a[0])
	f.When = X.ToU(a[1])
	return f
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (f *Foo) FindOffsetLimit(offset, limit uint32, idx string) []Foo { //nolint:dupl false positive
	var rows []Foo
	res, err := f.Adapter.Connection.Do(
		tarantool.NewSelectRequest(f.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).Get()
	if L.IsError(err, `Foo.FindOffsetLimit failed: `+f.SpaceName()) {
		return rows
	}
	for _, row := range res {
		item := Foo{}
		row, ok := row.([]any)
		if ok {
			rows = append(rows, *item.FromArray(row))
		}
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (f *Foo) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	resp, err := f.Adapter.Connection.Do(
		tarantool.NewSelectRequest(f.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).GetResponse()
	if L.IsError(err, `Foo.FindOffsetLimit failed: `+f.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	res, err := resp.Decode()
	if L.IsError(err, `Foo.FindOffsetLimit failed: `+f.SpaceName()) {
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
func (f *Foo) Total() int64 { //nolint:dupl false positive
	rows := f.Adapter.CallBoxSpace(f.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// FooFieldTypeMap returns key value of field name and key
var FooFieldTypeMap = map[string]Tt.DataType { //nolint:dupl false positive
	`id`:   Tt.Unsigned,
	`when`: Tt.Unsigned,
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

