package rqPoints

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	`hugedbbench/2023geo/tarantool/mPoints`

	`github.com/tarantool/go-tarantool/v2`

	`github.com/kokizzu/gotro/A`
	`github.com/kokizzu/gotro/D/Tt`
	`github.com/kokizzu/gotro/L`
	`github.com/kokizzu/gotro/X`
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type rqPoints__ORM.GEN.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type rqPoints__ORM.GEN.go
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
	return string(mPoints.TablePointsSg) // casting required to string from Tt.TableName
}

// SqlTableName returns quoted table name
func (p *PointsSg) SqlTableName() string { //nolint:dupl false positive
	return `"points_sg"`
}

func (p *PointsSg) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

// FindById Find one by Id
func (p *PointsSg) FindById() bool { //nolint:dupl false positive
	res, err := p.Adapter.Connection.Do(
		tarantool.NewSelectRequest(p.SpaceName()).
		Index(p.UniqueIndexId()).
		Offset(0).
		Limit(1).
		Iterator(tarantool.IterEq).
		Key(A.X{p.Id}),
	).Get()
	if L.IsError(err, `PointsSg.FindById failed: `+p.SpaceName()) {
		return false
	}
	if len(res) == 1 {
		if row, ok := res[0].([]any); ok {
			p.FromArray(row)
			return true
		}
	}
	return false
}

// SpatialIndexCoord return spatial index name
func (p *PointsSg) SpatialIndexCoord() string { //nolint:dupl false positive
	return `coord`
}

// SqlSelectAllFields generate Sql select fields
func (p *PointsSg) SqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "coord"
	`
}

// SqlSelectAllUncensoredFields generate Sql select fields
func (p *PointsSg) SqlSelectAllUncensoredFields() string { //nolint:dupl false positive
	return ` "id"
	, "coord"
	`
}

// ToUpdateArray generate slice of update command
func (p *PointsSg) ToUpdateArray() *tarantool.Operations { //nolint:dupl false positive
	return tarantool.NewOperations().
		Assign(0, p.Id).
		Assign(1, p.Coord)
}

// IdxId return name of the index
func (p *PointsSg) IdxId() int { //nolint:dupl false positive
	return 0
}

// SqlId return name of the column being indexed
func (p *PointsSg) SqlId() string { //nolint:dupl false positive
	return `"id"`
}

// IdxCoord return name of the index
func (p *PointsSg) IdxCoord() int { //nolint:dupl false positive
	return 1
}

// SqlCoord return name of the column being indexed
func (p *PointsSg) SqlCoord() string { //nolint:dupl false positive
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

// FromUncensoredArray convert slice to receiver fields
func (p *PointsSg) FromUncensoredArray(a A.X) *PointsSg { //nolint:dupl false positive
	p.Id = X.ToU(a[0])
	p.Coord = X.ToArr(a[1])
	return p
}

// FindOffsetLimit returns slice of struct, order by idx, eg. .UniqueIndex*()
func (p *PointsSg) FindOffsetLimit(offset, limit uint32, idx string) []PointsSg { //nolint:dupl false positive
	var rows []PointsSg
	res, err := p.Adapter.Connection.Do(
		tarantool.NewSelectRequest(p.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).Get()
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
		return rows
	}
	for _, row := range res {
		item := PointsSg{}
		row, ok := row.([]any)
		if ok {
			rows = append(rows, *item.FromArray(row))
		}
	}
	return rows
}

// FindArrOffsetLimit returns as slice of slice order by idx eg. .UniqueIndex*()
func (p *PointsSg) FindArrOffsetLimit(offset, limit uint32, idx string) ([]A.X, Tt.QueryMeta) { //nolint:dupl false positive
	var rows []A.X
	resp, err := p.Adapter.Connection.Do(
		tarantool.NewSelectRequest(p.SpaceName()).
		Index(idx).
		Offset(offset).
		Limit(limit).
		Iterator(tarantool.IterAll).
		Key(A.X{}),
	).GetResponse()
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
		return rows, Tt.QueryMetaFrom(resp, err)
	}
	res, err := resp.Decode()
	if L.IsError(err, `PointsSg.FindOffsetLimit failed: `+p.SpaceName()) {
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
func (p *PointsSg) Total() int64 { //nolint:dupl false positive
	rows := p.Adapter.CallBoxSpace(p.SpaceName() + `:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

// PointsSgFieldTypeMap returns key value of field name and key
var PointsSgFieldTypeMap = map[string]Tt.DataType { //nolint:dupl false positive
	`id`:    Tt.Unsigned,
	`coord`: Tt.Array,
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

