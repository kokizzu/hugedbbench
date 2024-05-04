package wcPoints

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2023geo/tarantool/mPoints/rqPoints`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/M`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type wcPoints__ORM.GEN.go
// PointsSgMutator DAO writer/command struct
type PointsSgMutator struct {
	rqPoints.PointsSg
	mutations *tarantool.Operations
	logs	  []A.X
}

// NewPointsSgMutator create new ORM writer/command object
func NewPointsSgMutator(adapter *Tt.Adapter) (res *PointsSgMutator) {
	res = &PointsSgMutator{PointsSg: rqPoints.PointsSg{Adapter: adapter}}
	res.mutations = tarantool.NewOperations()
	res.Coord = []any{}
	return
}

// Logs get array of logs [field, old, new]
func (p *PointsSgMutator) Logs() []A.X { //nolint:dupl false positive
	return p.logs
}

// HaveMutation check whether Set* methods ever called
func (p *PointsSgMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(p.logs) > 0
}

// ClearMutations clear all previously called Set* methods
func (p *PointsSgMutator) ClearMutations() { //nolint:dupl false positive
	p.mutations = tarantool.NewOperations()
	p.logs = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (p *PointsSgMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Connection.Do(tarantool.NewUpdateRequest(p.SpaceName()).
		Index(p.UniqueIndexId()).
		Key(A.X{p.Id}).
		Operations(p.ToUpdateArray()),
	).Get()
	return !L.IsError(err, `PointsSg.DoOverwriteById failed: `+p.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *PointsSgMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Connection.Do(
		tarantool.NewUpdateRequest(p.SpaceName()).
		Index(p.UniqueIndexId()).
		Key(A.X{p.Id}).
		Operations(p.mutations),
	).Get()
	return !L.IsError(err, `PointsSg.DoUpdateById failed: `+p.SpaceName())
}

// DoDeletePermanentById permanent delete
func (p *PointsSgMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Connection.Do(
		tarantool.NewDeleteRequest(p.SpaceName()).
		Index(p.UniqueIndexId()).
		Key(A.X{p.Id}),
	).Get()
	return !L.IsError(err, `PointsSg.DoDeletePermanentById failed: `+p.SpaceName())
}

// func (p *PointsSgMutator) DoUpsert() bool { //nolint:dupl false positive
//	arr := p.ToArray()
//	_, err := p.Adapter.Upsert(p.SpaceName(), arr, A.X{
//		A.X{`=`, 0, p.Id},
//		A.X{`=`, 1, p.Coord},
//	})
//	return !L.IsError(err, `PointsSg.DoUpsert failed: `+p.SpaceName()+ `\n%#v`, arr)
// }

// DoInsert insert, error if already exists
func (p *PointsSgMutator) DoInsert() bool { //nolint:dupl false positive
	arr := p.ToArray()
	row, err := p.Adapter.Connection.Do(
		tarantool.NewInsertRequest(p.SpaceName()).
		Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				p.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `PointsSg.DoInsert failed: `+p.SpaceName() + `\n%#v`, arr)
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (p *PointsSgMutator) DoUpsert() bool { //nolint:dupl false positive
	arr := p.ToArray()
	row, err := p.Adapter.Connection.Do(
		tarantool.NewReplaceRequest(p.SpaceName()).
		Tuple(arr),
	).Get()
	if err == nil {
		if len(row) > 0 {
			if cells, ok := row[0].([]any); ok && len(cells) > 0 {
				p.Id = X.ToU(cells[0])
			}
		}
	}
	return !L.IsError(err, `PointsSg.DoUpsert failed: `+p.SpaceName()+ `\n%#v`, arr)
}

// SetId create mutations, should not duplicate
func (p *PointsSgMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != p.Id {
		p.mutations.Assign(0, val)
		p.logs = append(p.logs, A.X{`id`, p.Id, val})
		p.Id = val
		return true
	}
	return false
}

// SetCoord create mutations, should not duplicate
func (p *PointsSgMutator) SetCoord(val []any) bool { //nolint:dupl false positive
	p.mutations.Assign(1, val)
	p.logs = append(p.logs, A.X{`coord`, p.Coord, val})
	p.Coord = val
	return true
}

// SetAll set all from another source, only if another property is not empty/nil/zero or in forceMap
func (p *PointsSgMutator) SetAll(from rqPoints.PointsSg, excludeMap, forceMap M.SB) (changed bool) { //nolint:dupl false positive
	if excludeMap == nil { // list of fields to exclude
		excludeMap = M.SB{}
	}
	if forceMap == nil { // list of fields to force overwrite
		forceMap = M.SB{}
	}
	if !excludeMap[`id`] && (forceMap[`id`] || from.Id != 0) {
		p.Id = from.Id
		changed = true
	}
	if !excludeMap[`coord`] && (forceMap[`coord`] || from.Coord != nil) {
		p.Coord = from.Coord
		changed = true
	}
	return
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

