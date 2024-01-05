package loopVsWhereIn

import (
	"testing"

	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/sourcegraph/conc/pool"
	"github.com/zeebo/assert"

	"loopVsWhereIn/mTest"
	"loopVsWhereIn/mTest/rqTest"
	"loopVsWhereIn/mTest/wcTest"
)

var taran *Tt.Adapter

const queryAll = `SELECT * FROM "test_table2"`
const queryOne = `SELECT * FROM "test_table2" WHERE "content" = `
const limit1k = ` LIMIT 1000`

func BenchmarkInsertS_Taran_ORM(b *testing.B) {
	if done() {
		b.SkipNow()
		return
	}
	defer timing()()
	b.N = total
	r := taran.ExecSql(`DELETE FROM ` + S.ZZ(mTest.TableTestTable2))
	assert.Equal(b, len(r), 1)

	p := pool.New().WithMaxGoroutines(2)
	for z := uint64(1); z <= total; z++ {
		z := z
		p.Go(func() {
			row := wcTest.NewTestTable2Mutator(taran)
			row.Id = z
			row.Content = S.EncodeCB63(z, 0)
			ok := row.DoUpsert()
			assert.True(b, ok)
		})
	}
	p.Wait()
}
func BenchmarkUpdate_Taran_ORM(b *testing.B) {
	if done() {
		b.SkipNow()
		return
	}
	defer timing()(2)
	b.N = total

	p := pool.New().WithMaxGoroutines(2)
	for z := uint64(1); z <= total; z++ {
		z := z
		p.Go(func() {
			row := wcTest.NewTestTable2Mutator(taran)
			row.Id = z
			row.Content = S.EncodeCB63(total+z, 0)
			ok := row.DoUpdateById()
			assert.True(b, ok)
			row = wcTest.NewTestTable2Mutator(taran) // create new mutator just to be fair
			row.Id = z
			row.Content = S.EncodeCB63(z, 0)
			ok = row.DoUpdateById()
			assert.True(b, ok)
		})
	}
	b.N *= 2
	p.Wait()
}

func BenchmarkGetAllStruct_Taran_SQL(b *testing.B) {

	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			res := make([]*mTest.TestTable2, 0, total)
			_ = taran.QuerySql(queryAll+limit1k, func(row []any) {
				obj := &mTest.TestTable2{}
				obj.FromArray(row)
				res = append(res, obj)
			})
			assert.Equal(b, len(res), limit)
		})
	}
	p.Wait()
}
func BenchmarkGetAllStruct_Taran_ORM(b *testing.B) {

	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			row := rqTest.NewTestTable2(taran)
			rows := row.FindOffsetLimit(0, limit, Tt.IdCol)
			assert.Equal(b, len(rows), limit)
		})
	}
	p.Wait()
}

func BenchmarkGetAllArray_Taran_ORM(b *testing.B) {

	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			row := rqTest.NewTestTable2(taran)
			rows, meta := row.FindArrOffsetLimit(0, limit, Tt.IdCol)
			if len(rows) != limit {
				L.Describe(meta)
				assert.Equal(b, len(rows), limit)
			}
		})
	}
	p.Wait()
}
func BenchmarkGetAllMap_Taran_SQL(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			obj := &mTest.TestTable2{}
			res := make([]map[string]any, 0, total)
			_ = taran.QuerySql(queryAll+limit1k, func(row []any) {
				m := obj.ToMapFromSlice(row)
				res = append(res, m)
			})
		})
	}
	p.Wait()
}

func BenchmarkGetOneStruct_Taran_SQL(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			obj := mTest.TestTable2{}
			_ = taran.QuerySql(queryOne+S.Z(
				S.EncodeCB63(i, 0)), func(row []any) {
				obj.FromArray(row)
			})
		})
	}
	p.Wait()
}

func BenchmarkGetOneMap_Taran_SQL(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			obj := &mTest.TestTable2{}
			_ = taran.QuerySql(queryOne+S.Z(
				S.EncodeCB63(i, 0)), func(row []any) {
				_ = obj.ToMapFromSlice(row)
			})
		})
	}
	p.Wait()
}

func BenchmarkGetOneStruct_Taran_ORM(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			obj := rqTest.NewTestTable2(taran)
			obj.Content = S.EncodeCB63(1+i%total, 0)
			ok := obj.FindByContent()
			assert.True(b, ok)
		})
	}
	p.Wait()
}

func BenchmarkGetWhereInStruct_Taran_SQL(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			obj := rqTest.NewTestTable2(taran)
			rows := obj.FindWhereInStruct(idsToFetch(i))
			assert.True(b, len(rows) > 0)
		})
	}
	p.Wait()
}

func BenchmarkGetWhereInArray_Taran_SQL(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			obj := rqTest.NewTestTable2(taran)
			rows := obj.FindWhereInArray(idsToFetch(i))
			assert.True(b, len(rows) > 0)
		})
	}
	p.Wait()
}

func BenchmarkGetLoopStruct_Taran_ORM(b *testing.B) {
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		i := i
		p.Go(func() {
			ids := idsToFetch(i)
			obj := rqTest.NewTestTable2(taran)
			res := make([]mTest.TestTable2, 0, 4)
			for _, id := range ids {
				copy := mTest.TestTable2{}
				obj.Id = id
				if obj.FindById() { // simulate append
					copy.Id = obj.Id
					copy.Content = obj.Content
					res = append(res, copy)
				}
			}
			assert.True(b, len(res) > 0)
		})
	}
	p.Wait()
}
