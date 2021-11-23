package wcFoo

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2021mq/tarantoolAsQ/mFoo/rqFoo`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcFoo__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type wcFoo__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type wcFoo__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type wcFoo__ORM.GEN.go
// go:generate msgp -tests=false -file wcFoo__ORM.GEN.go -o wcFoo__MSG.GEN.go

type FooMutator struct {
	rqFoo.Foo
	mutations []A.X
}

func NewFooMutator(adapter *Tt.Adapter) *FooMutator {
	return &FooMutator{Foo: rqFoo.Foo{Adapter: adapter}}
}

func (f *FooMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(f.mutations) > 0
}

// Overwrite all columns, error if not exists
func (f *FooMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := f.Adapter.Update(f.SpaceName(), f.UniqueIndexId(), A.X{f.Id}, f.ToUpdateArray())
	return !L.IsError(err, `Foo.DoOverwriteById failed: `+f.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (f *FooMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !f.HaveMutation() {
		return true
	}
	_, err := f.Adapter.Update(f.SpaceName(), f.UniqueIndexId(), A.X{f.Id}, f.mutations)
	return !L.IsError(err, `Foo.DoUpdateById failed: `+f.SpaceName())
}

func (f *FooMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := f.Adapter.Delete(f.SpaceName(), f.UniqueIndexId(), A.X{f.Id})
	return !L.IsError(err, `Foo.DoDeletePermanentById failed: `+f.SpaceName())
}

// func (f *FooMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := f.Adapter.Upsert(f.SpaceName(), f.ToArray(), A.X{
//		A.X{`=`, 0, f.Id},
//		A.X{`=`, 1, f.When},
//	})
//	return !L.IsError(err, `Foo.DoUpsert failed: `+f.SpaceName())
// }

// insert, error if exists
func (f *FooMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := f.Adapter.Insert(f.SpaceName(), f.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			f.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Foo.DoInsert failed: `+f.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (f *FooMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := f.Adapter.Replace(f.SpaceName(), f.ToArray())
	return !L.IsError(err, `Foo.DoReplace failed: `+f.SpaceName())
}

func (f *FooMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != f.Id {
		f.mutations = append(f.mutations, A.X{`=`, 0, val})
		f.Id = val
		return true
	}
	return false
}

func (f *FooMutator) SetWhen(val uint64) bool { //nolint:dupl false positive
	if val != f.When {
		f.mutations = append(f.mutations, A.X{`=`, 1, val})
		f.When = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

