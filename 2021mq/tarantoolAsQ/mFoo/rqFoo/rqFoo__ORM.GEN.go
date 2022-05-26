package rqFoo

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"hugedbbench/2021mq/tarantoolAsQ/mFoo"

	"github.com/tarantool/go-tarantool"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqFoo__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type rqFoo__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type rqFoo__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type rqFoo__ORM.GEN.go
// go:generate msgp -tests=false -file rqFoo__ORM.GEN.go -o rqFoo__MSG.GEN.go

type Foo struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id      uint64
	When    uint64
}

func NewFoo(adapter *Tt.Adapter) *Foo {
	return &Foo{Adapter: adapter}
}

func (f *Foo) SpaceName() string { //nolint:dupl false positive
	return string(mFoo.TableFoo)
}

func (f *Foo) sqlTableName() string { //nolint:dupl false positive
	return `"foo"`
}

func (f *Foo) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

func (f *Foo) FindById() bool { //nolint:dupl false positive
	res, err := f.Adapter.Select(f.SpaceName(), f.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{f.Id})
	if L.IsError(err, `Foo.FindById failed: `+f.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		f.FromArray(rows[0])
		return true
	}
	return false
}

func (f *Foo) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "when"
	`
}

func (f *Foo) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, f.Id},
		A.X{`=`, 1, f.When},
	}
}

func (f *Foo) IdxId() int { //nolint:dupl false positive
	return 0
}

func (f *Foo) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (f *Foo) IdxWhen() int { //nolint:dupl false positive
	return 1
}

func (f *Foo) sqlWhen() string { //nolint:dupl false positive
	return `"when"`
}

func (f *Foo) ToArray() A.X { //nolint:dupl false positive
	var id interface{} = nil
	if f.Id != 0 {
		id = f.Id
	}
	return A.X{
		id,
		f.When, // 1
	}
}

func (f *Foo) FromArray(a A.X) *Foo { //nolint:dupl false positive
	f.Id = X.ToU(a[0])
	f.When = X.ToU(a[1])
	return f
}

func (f *Foo) Total() int64 { //nolint:dupl false positive
	rows := f.Adapter.CallBoxSpace(f.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
