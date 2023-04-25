package wcPoints

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"hugedbbench/2023geo/tarantool/mPoints/rqPoints"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

// PointsSgMutator DAO writer/command struct
//
//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'Id" form' 'Id,string" form' type wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'json:"id"' 'json:"id,string"' type wcPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'By" form' 'By,string" form' type wcPoints__ORM.GEN.go
type PointsSgMutator struct {
	rqPoints.PointsSg
	mutations []A.X
}

// NewPointsSgMutator create new ORM writer/command object
func NewPointsSgMutator(adapter *Tt.Adapter) *PointsSgMutator {
	return &PointsSgMutator{PointsSg: rqPoints.PointsSg{Adapter: adapter}}
}

// HaveMutation check whether Set* methods ever called
func (p *PointsSgMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(p.mutations) > 0
}

// ClearMutations clear all previously called Set* methods
func (p *PointsSgMutator) ClearMutations() { //nolint:dupl false positive
	p.mutations = []A.X{}
}

// DoOverwriteById update all columns, error if not exists, not using mutations/Set*
func (p *PointsSgMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.ToUpdateArray())
	return !L.IsError(err, `PointsSg.DoOverwriteById failed: `+p.SpaceName())
}

// DoUpdateById update only mutated fields, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *PointsSgMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.mutations)
	return !L.IsError(err, `PointsSg.DoUpdateById failed: `+p.SpaceName())
}

// DoDeletePermanentById permanent delete
func (p *PointsSgMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Delete(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id})
	return !L.IsError(err, `PointsSg.DoDeletePermanentById failed: `+p.SpaceName())
}

// func (p *PointsSgMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := p.Adapter.Upsert(p.SpaceName(), p.ToArray(), A.X{
//		A.X{`=`, 0, p.Id},
//		A.X{`=`, 1, p.Coord},
//	})
//	return !L.IsError(err, `PointsSg.DoUpsert failed: `+p.SpaceName())
// }

// DoInsert insert, error if already exists
func (p *PointsSgMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := p.Adapter.Insert(p.SpaceName(), p.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			p.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `PointsSg.DoInsert failed: `+p.SpaceName())
}

// DoUpsert upsert, insert or overwrite, will error only when there's unique secondary key being violated
// replace = upsert, only error when there's unique secondary key
// previous name: DoReplace
func (p *PointsSgMutator) DoUpsert() bool { //nolint:dupl false positive
	_, err := p.Adapter.Replace(p.SpaceName(), p.ToArray())
	return !L.IsError(err, `PointsSg.DoUpsert failed: `+p.SpaceName())
}

// SetId create mutations, should not duplicate
func (p *PointsSgMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != p.Id {
		p.mutations = append(p.mutations, A.X{`=`, 0, val})
		p.Id = val
		return true
	}
	return false
}

// SetCoord create mutations, should not duplicate
func (p *PointsSgMutator) SetCoord(val []any) bool { //nolint:dupl false positive
	p.mutations = append(p.mutations, A.X{`=`, 1, val})
	p.Coord = val
	return true
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
