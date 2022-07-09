/*
 * Copyright (c) 2022, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL
 * WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL
 * THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR
 * CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING
 * FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF
 * CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package gendb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/results"
)

// Exec calls ExecContext with the default background context.
func Exec(
	db *sql.DB,
	query string,
	args ...any,
) *genfuncs.Result[sql.Result] {
	return ExecContext(db, context.Background(), query, args...)
}

// ExecContext executes a query with arguments and returns a sql.Result summarizing the effect of the statement.
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

// Query calls QueryContext with the default background context.
func Query[T any](
	db *sql.DB,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[[]T] {
	return QueryContext[T](db, context.Background(), binder, query, args...)
}

// QueryContext  performs a query with arguments using a binder to assign the rows returned to a slice of results.
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

// QueryRow calls QueryRowContext with the default Background context.
func QueryRow[T any](
	db *sql.DB,
	binder func(*T) []any,
	query string,
	args ...any,
) *genfuncs.Result[T] {
	return QueryRowContext[T](db, context.Background(), binder, query, args...)
}

// QueryRowContext performs a query with arguments using a binder to assign the first row returned to a single result.
// If multiple rows are returned the first one is used. If no rows are returned the error sql.ErrNoRows is returned.
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

// SingleBinder is a binder for base types.
func SingleBinder[T any](t *T) []any {
	return []any{t}
}
