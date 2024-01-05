package loopVsWhereIn

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kokizzu/gotro/S"
	"github.com/sourcegraph/conc/pool"
	"github.com/stretchr/testify/assert"
)

func benchmarkInsertPgx(b *testing.B, pgxConn *pgxpool.Pool) {
	if done() {
		b.SkipNow()
		return
	}
	defer timing()()
	b.N = total
	ctx := context.Background()
	_, err := pgxConn.Exec(ctx, `TRUNCATE TABLE `+pgxTableName)
	assert.Nil(b, err)
	p := pool.New().WithMaxGoroutines(cores)
	for z := uint64(1); z <= total; z++ {
		z := z
		p.Go(func() {
			_, err := pgxConn.Exec(ctx, `INSERT INTO `+pgxTableName+` (id, content) VALUES ($1, $2)`, z, S.EncodeCB63(z, 0))
			assert.Nil(b, err)
		})
	}
	p.Wait()
}

func benchmarkUpdatePgx(b *testing.B, pgxConn *pgxpool.Pool) {
	if done() {
		b.SkipNow()
		return
	}
	defer timing()(2)
	b.N = total
	ctx := context.Background()
	p := pool.New().WithMaxGoroutines(cores)
	for z := uint64(1); z <= total; z++ {
		z := z
		p.Go(func() {
			_, err := pgxConn.Exec(ctx, `UPDATE  `+pgxTableName+` SET content = $2 WHERE id = $1`, z, S.EncodeCB63(total+z, 0))
			assert.Nil(b, err)
			_, err = pgxConn.Exec(ctx, `UPDATE  `+pgxTableName+` SET content = $2 WHERE id = $1`, z, S.EncodeCB63(z, 0))
			assert.Nil(b, err)
		})
	}
	b.N *= 2
	p.Wait()
}

func benchmarkGetOneStructPgx(b *testing.B, pgxConn *pgxpool.Pool) {
	ctx := context.Background()
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			row := pgxConn.QueryRow(ctx, `SELECT * FROM `+pgxTableName+` WHERE content = $1 LIMIT 1`, S.EncodeCB63(1+(i%total), 0))
			var row2 PgxTestTable
			err := row.Scan(&row2.Id, &row2.Content)
			assert.Nil(b, err)
		})
	}
	p.Wait()
}
func benchmarkGetAllStructPgx(b *testing.B, pgxConn *pgxpool.Pool) {
	ctx := context.Background()
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			rows, err := pgxConn.Query(ctx, `SELECT * FROM `+pgxTableName+` LIMIT $1`, limit)
			assert.Nil(b, err)
			defer rows.Close()
			res := make([]PgxTestTable, 0, limit)
			for rows.Next() {
				var row PgxTestTable
				err = rows.Scan(&row.Id, &row.Content)
				assert.Nil(b, err)
				res = append(res, row)
			}
			assert.Equal(b, len(res), limit)
		})
	}
	p.Wait()
}

func benchmarkWhereInPgx(b *testing.B, pgxConn *pgxpool.Pool) {
	ctx := context.Background()
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			it, err := pgxConn.Query(ctx, `SELECT * FROM `+pgxTableName+` WHERE id IN($1,$2,$3,$4) LIMIT 1`, idsToFetchAny(i)...)
			assert.Nil(b, err)
			var row2 PgxTestTable
			rows := make([]PgxTestTable, 0, 4)
			for it.Next() {
				err := it.Scan(&row2.Id, &row2.Content)
				if errors.Is(err, pgx.ErrNoRows) {
					break
				}
				assert.Nil(b, err)
				rows = append(rows, row2)
			}
			assert.True(b, len(rows) > 0)
			it.Close()
		})
	}
	p.Wait()
}

func benchmarkLoopPgx(b *testing.B, pgxConn *pgxpool.Pool) {
	ctx := context.Background()
	p := pool.New().WithMaxGoroutines(cores)
	for i := uint64(1); i <= uint64(b.N); i++ {
		p.Go(func() {
			ids := idsToFetch(i)
			rows := make([]PgxTestTable, 0, 4)
			for _, id := range ids {
				row := pgxConn.QueryRow(ctx, `SELECT * FROM `+pgxTableName+` WHERE id = $i LIMIT 1`, id)
				var row2 PgxTestTable
				err := row.Scan(&row2.Id, &row2.Content)
				_ = errors.Is(err, pgx.ErrNoRows)
				rows = append(rows, row2)
			}
			assert.True(b, len(rows) > 0)
		})
	}
	p.Wait()
}
