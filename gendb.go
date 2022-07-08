package gendb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/results"
)

func Exec(
	db *sql.DB,
	query string,
	args ...any,
) *genfuncs.Result[sql.Result] {
	return ExecContext(db, context.Background(), query, args...)
}

func ExecContext(
	db *sql.DB,
	ctx context.Context,
	query string,
	args ...any,
) *genfuncs.Result[sql.Result] {
	return results.Map[*sql.Stmt, sql.Result](
		genfuncs.NewResultError(db.Prepare(query)),
		func(statement *sql.Stmt) (result *genfuncs.Result[sql.Result]) {
			defer func() {
				_ = statement.Close()
			}()
			return genfuncs.NewResultError(statement.ExecContext(ctx, args...))
		})
}

func Query[T any](
	db *sql.DB,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[[]T] {
	return QueryContext[T](db, context.Background(), binder, query, args...)
}

func QueryContext[T any](
	db *sql.DB,
	ctx context.Context,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[[]T] {
	return results.Map[*sql.Rows, []T](
		genfuncs.NewResultError(db.QueryContext(ctx, query, args...)),
		func(rows *sql.Rows) *genfuncs.Result[[]T] {
			defer func() {
				_ = rows.Close()
			}()
			var values []T
			for rows.Next() {
				var t T
				cols := binder(&t)
				err := rows.Scan(cols...)
				if err != nil {
					return genfuncs.NewError[[]T](fmt.Errorf("SQL: %s\n\tscan failed %w", query, err))
				}
				values = append(values, t)
			}
			return genfuncs.NewResult(values)
		})
}

func QueryRow[T any](
	db *sql.DB,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[T] {
	return QueryRowContext[T](db, context.Background(), binder, query, args...)
}

func QueryRowContext[T any](
	db *sql.DB,
	ctx context.Context,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[T] {
	return results.Map[*sql.Rows, T](
		genfuncs.NewResultError(db.QueryContext(ctx, query, args...)),
		func(rows *sql.Rows) *genfuncs.Result[T] {
			defer func() {
				_ = rows.Close()
			}()
			if !rows.Next() {
				return genfuncs.NewError[T](sql.ErrNoRows)
			}
			var t T
			cols := binder(&t)
			err := rows.Scan(cols...)
			if err != nil {
				return genfuncs.NewError[T](fmt.Errorf("SQL: %s\n\tscan failed %w", query, err))
			}
			return genfuncs.NewResult(t)
		})
}

func SingleBinder[T any](t *T) []any {
	return []any{t}
}
