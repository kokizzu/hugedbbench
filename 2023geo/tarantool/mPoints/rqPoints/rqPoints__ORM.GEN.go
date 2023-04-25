package rqPoints

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2023geo/tarantool/mPoints`

	`github.com/tarantool/go-tarantool`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'Id" form' 'Id,string" form' type rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'json:"id"' 'json:"id,string"' type rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix 'By" form' 'By,string" form' type rqPoints__ORM.GEN.go
// go:generate msgp -tests=false -file rqPoints__ORM.GEN.go -o rqPoints__MSG.GEN.go

// PointsSg DAO reader/query struct
type PointsSg struct {
	Adapter *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-"`
	Id    uint64
	Coord []any
}

// NewPointsSg create new ORM reader/query object
func NewPointsSg(adapter *Tt.Adapter) *PointsSg {
	return &PointsSg{Adapter: adapter}
}

// SpaceName returns full package and table name
func (p *PointsSg) SpaceName() string { //nolint:dupl false positive
	return string(mPoints.TablePointsSg)
}

// sqlTableName returns quoted table name
func (p *PointsSg) sqlTableName() string { //nolint:dupl false positive
	return `"points_sg"`
}

func (p *PointsSg) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (p *PointsSg) FindById() bool { //nolint:dupl false positive
	res, err := p.Adapter.Select(p.SpaceName(), p.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{p.Id})
	if L.IsError(err, `PointsSg.FindById failed: `+p.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		p.FromArray(rows[0])
		return true
	}
	return false
}

// SpatialIndexCoord return spatial index name
func (p *PointsSg) SpatialIndexCoord() string { //nolint:dupl false positive
	return `coord`
}

// sqlSelectAllFields generate sql select fields
func (p *PointsSg) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "coord"
	`
}

// ToUpdateArray generate slice of update command
func (p *PointsSg) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, p.Id},
		A.X{`=`, 1, p.Coord},
	}
}

// IdxId return name of the index
func (p *PointsSg) IdxId() int { //nolint:dupl false positive
	return 0
}

// sqlId return name of the column being indexed
func (p *PointsSg) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxCoord return name of the index
func (p *PointsSg) IdxCoord() int { //nolint:dupl false positive
	return 1
}

// sqlCoord return name of the column being indexed
func (p *PointsSg) sqlCoord() string { //nolint:dupl false positive
	return `"coord"`
}

// ToArray receiver fields to slice
func (p *PointsSg) ToArray() A.X { //nolint:dupl false positive
	var id any = nil
	if p.Id != 0 {
		id = p.Id
	}
	return A.X{
		id,
		p.Coord, // 1
	}
}

// FromArray convert slice to receiver fields
func (p *PointsSg) FromArray(a A.X) *PointsSg { //nolint:dupl false positive
	p.Id = X.ToU(a[0])
	p.Coord = X.ToArr(a[1])
	return p
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (p *PointsSg) FindOffsetLimit(offset, limit uint32, idx string) []PointsSg { //nolint:dupl false positive
	var rows []PointsSg
	res, err := p.Adapter.Select(p.SpaceName(), idx, offset, limit, 2, A.X{})
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
		return rows
	}
	for _, row := range res.Tuples() {
		item := PointsSg{}
		rows = append(rows, *item.FromArray(row))
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (p *PointsSg) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	res, err := p.Adapter.Select(p.SpaceName(), idx, offset, limit, 2, A.X{})
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
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
func (p *PointsSg) Total() int64 { //nolint:dupl false positive
	rows := p.Adapter.CallBoxSpace(p.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

