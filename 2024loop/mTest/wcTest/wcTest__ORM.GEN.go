package wcTest

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"loopVsWhereIn/mTest/rqTest"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
)

// TestTable2Mutator DAO writer/command struct
//
//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcTest__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcTest__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcTest__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcTest__ORM.GEN.go
type TestTable2Mutator struct {
	rqTest.TestTable2
	mutations []A.X
	logs      []A.X
}

// NewTestTable2Mutator create new ORM writer/command object
func NewTestTable2Mutator(adapter *Tt.Adapter) (res *TestTable2Mutator) {
	res = &TestTable2Mutator{TestTable2: rqTest.TestTable2{Adapter: adapter}}
	return
}

// Logs get array of logs [field, old, new]
func (t *TestTable2Mutator) Logs() []A.X { //nolint:dupl false positive
	return t.logs
}

// HaveMutation check whether Set* methods ever called
func (t *TestTable2Mutator) HaveMutation() bool { //nolint:dupl false positive
	return len(t.mutations) > 0
}

// ClearMutations clear all previously called Set* methods
func (t *TestTable2Mutator) ClearMutations() { //nolint:dupl false positive
	t.mutations = []A.X{}
	t.logs = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (t *TestTable2Mutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := t.Adapter.Update(t.SpaceName(), t.UniqueIndexId(), A.X{t.Id}, t.ToUpdateArray())
	return !L.IsError(err, `TestTable2.DoOverwriteById failed: `+t.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (t *TestTable2Mutator) DoUpdateById() bool { //nolint:dupl false positive
	if !t.HaveMutation() {
		return true
	}
	_, err := t.Adapter.Update(t.SpaceName(), t.UniqueIndexId(), A.X{t.Id}, t.mutations)
	return !L.IsError(err, `TestTable2.DoUpdateById failed: `+t.SpaceName())
}

// DoDeletePermanentById permanent delete
func (t *TestTable2Mutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := t.Adapter.Delete(t.SpaceName(), t.UniqueIndexId(), A.X{t.Id})
	return !L.IsError(err, `TestTable2.DoDeletePermanentById failed: `+t.SpaceName())
}

// func (t *TestTable2Mutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := t.ToArray()
//	_, err := t.Adapter.Upsert(t.SpaceName(), arr, A.X{
//		A.X{`=`, 0, t.Id},
//		A.X{`=`, 1, t.Content},
//	})
//	return !L.IsError(err, `TestTable2.DoUpsert failed: `+t.SpaceName()+ `\n%#v`, arr)
// }

// DoOverwriteByContent update all columns, error if not exists, not using mutations/Set*
func (t *TestTable2Mutator) DoOverwriteByContent() bool { //nolint:dupl false positive
	_, err := t.Adapter.Update(t.SpaceName(), t.UniqueIndexContent(), A.X{t.Content}, t.ToUpdateArray())
	return !L.IsError(err, `TestTable2.DoOverwriteByContent failed: `+t.SpaceName())
}

// DoUpdateByContent update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (t *TestTable2Mutator) DoUpdateByContent() bool { //nolint:dupl false positive
	if !t.HaveMutation() {
		return true
	}
	_, err := t.Adapter.Update(t.SpaceName(), t.UniqueIndexContent(), A.X{t.Content}, t.mutations)
	return !L.IsError(err, `TestTable2.DoUpdateByContent failed: `+t.SpaceName())
}

// DoDeletePermanentByContent permanent delete
func (t *TestTable2Mutator) DoDeletePermanentByContent() bool { //nolint:dupl false positive
	_, err := t.Adapter.Delete(t.SpaceName(), t.UniqueIndexContent(), A.X{t.Content})
	return !L.IsError(err, `TestTable2.DoDeletePermanentByContent failed: `+t.SpaceName())
}

// DoInsert insert, error if already exists
func (t *TestTable2Mutator) DoInsert() bool { //nolint:dupl false positive
	arr := t.ToArray()
	row, err := t.Adapter.Insert(t.SpaceName(), arr)
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			t.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `TestTable2.DoInsert failed: `+t.SpaceName()+`\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (t *TestTable2Mutator) DoUpsert() bool { //nolint:dupl false positive
	arr := t.ToArray()
	row, err := t.Adapter.Replace(t.SpaceName(), arr)
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			t.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `TestTable2.DoUpsert failed: `+t.SpaceName()+`\n%#v`, arr)
}

// SetId create mutations, should not duplicate
func (t *TestTable2Mutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != t.Id {
		t.mutations = append(t.mutations, A.X{`=`, 0, val})
		t.logs = append(t.logs, A.X{`id`, t.Id, val})
		t.Id = val
		return true
	}
	return false
}

// SetContent create mutations, should not duplicate
func (t *TestTable2Mutator) SetContent(val string) bool { //nolint:dupl false positive
	if val != t.Content {
		t.mutations = append(t.mutations, A.X{`=`, 1, val})
		t.logs = append(t.logs, A.X{`content`, t.Content, val})
		t.Content = val
		return true
	}
	return false
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (t *TestTable2Mutator) SetAll(from rqTest.TestTable2, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		t.Id = from.Id
		changed = true
	}
	if !excludeMap[`content`] && (forceMap[`content`] || from.Content != ``) {
		t.Content = S.Trim(from.Content)
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
