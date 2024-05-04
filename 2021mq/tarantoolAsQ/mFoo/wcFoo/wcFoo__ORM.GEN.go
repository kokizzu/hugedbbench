package wcFoo

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2021mq/tarantoolAsQ/mFoo/rqFoo`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/M`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcFoo__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcFoo__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcFoo__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcFoo__ORM.GEN.go
// FooMutator DAO writer/command struct
type FooMutator struct {
	rqFoo.Foo
	mutations *tarantool.Operations
	logs	  []A.X
}

// NewFooMutator create new ORM writer/command object
func NewFooMutator(adapter *Tt.Adapter) (res *FooMutator) {
	res = &FooMutator{Foo: rqFoo.Foo{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	return
}

// Logs get array of logs [field, old, new]
func (f *FooMutator) Logs() []A.X { //nolint:dupl false positive
	return f.logs
}

// HaveMutation check whether Set* methods ever called
func (f *FooMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(f.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (f *FooMutator) ClearMutations() { //nolint:dupl false positive
	f.mutations = tarantool.NewOperations()
	f.logs = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (f *FooMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := f.Adapter.Connection.Do(tarantool.NewUpdateRequest(f.SpaceName()).
		Index(f.UniqueIndexId()).
		Key(A.X{f.Id}).
		Operations(f.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `Foo.DoOverwriteById failed: `+f.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (f *FooMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !f.HaveMutation() {
		return true
	}
	_, err := f.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(f.SpaceName()).
		Index(f.UniqueIndexId()).
		Key(A.X{f.Id}).
		Operations(f.mutations),
	).Get()
	return !L.IsError(err, `Foo.DoUpdateById failed: `+f.SpaceName())
}

// DoDeletePermanentById permanent delete
func (f *FooMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := f.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(f.SpaceName()).
		Index(f.UniqueIndexId()).
		Key(A.X{f.Id}),
	).Get()
	return !L.IsError(err, `Foo.DoDeletePermanentById failed: `+f.SpaceName())
}

// func (f *FooMutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := f.ToArray()
//	_, err := f.Adapter.Upsert(f.SpaceName(), arr, A.X{
//		A.X{`=`, 0, f.Id},
//		A.X{`=`, 1, f.When},
//	})
//	return !L.IsError(err, `Foo.DoUpsert failed: `+f.SpaceName()+ `\n%#v`, arr)
// }

// DoInsert insert, error if already exists
func (f *FooMutator) DoInsert() bool { //nolint:dupl false positive
	arr := f.ToArray()
	row, err := f.Adapter.Connection.Do(
		tarantool.NewInsertRequest(f.SpaceName()).
		Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				f.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `Foo.DoInsert failed: `+f.SpaceName() + `\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (f *FooMutator) DoUpsert() bool { //nolint:dupl false positive
	arr := f.ToArray()
	row, err := f.Adapter.Connection.Do(
		tarantool.NewReplaceRequest(f.SpaceName()).
		Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				f.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `Foo.DoUpsert failed: `+f.SpaceName()+ `\n%#v`, arr)
}

// SetId create mutations, should not duplicate
func (f *FooMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != f.Id {
		f.mutations.Assign(0, val)
		f.logs = append(f.logs, A.X{`id`, f.Id, val})
		f.Id = val
		return true
	}
	return false
}

// SetWhen create mutations, should not duplicate
func (f *FooMutator) SetWhen(val uint64) bool { //nolint:dupl false positive
	if val != f.When {
		f.mutations.Assign(1, val)
		f.logs = append(f.logs, A.X{`when`, f.When, val})
		f.When = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (f *FooMutator) SetAll(from rqFoo.Foo, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		f.Id = from.Id
		changed = true
	}
	if !excludeMap[`when`] && (forceMap[`when`] || from.When != 0) {
		f.When = from.When
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

