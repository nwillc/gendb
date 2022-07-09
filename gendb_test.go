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

package gendb_test

import (
	"database/sql"
	"fmt"
	"gendb"
	"gendb/internal/sql_test"
	"github.com/nwillc/genfuncs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoodExec(t *testing.T) {
	db := sql_test.SetupDB(t)

	gendb.Exec(db, "INSERT INTO student(code, name, program) VALUES (?, ?, ?)", "0001", "fred", "masters").
		OnError(func(e error) {
			assert.NoError(t, e)
		}).
		OnSuccess(func(result sql.Result) {
			rows := genfuncs.NewResultError(result.RowsAffected())
			assert.True(t, rows.Ok())
			assert.Equal(t, int64(1), rows.OrEmpty())
		})
}

func TestBadExec(t *testing.T) {
	db := sql_test.SetupDB(t)

	gendb.Exec(db, "Foo").
		OnError(func(e error) {
			assert.Error(t, e)
		}).
		OnSuccess(func(result sql.Result) {
			assert.Fail(t, "expected Fail")
		})

	gendb.Exec(db, "INSERT INTO student(code, name, program) VALUES (?, ?, ?)", "0001", "fred").
		OnError(func(e error) {
			assert.Error(t, e)
		}).
		OnSuccess(func(result sql.Result) {
			assert.Fail(t, "expected Fail")
		})
}

func TestQuery_QueryRow(t *testing.T) {
	db := sql_test.SetupDB(t)

	// Test for an aggregation
	gendb.QueryRow[int](
		db,
		gendb.SingleBinder[int],
		"SELECT count(*) from STUDENT",
	).
		OnError(func(e error) {
			assert.NoError(t, e)
		}).
		OnSuccess(func(count int) {
			assert.Equal(t, 2, count)
		})

	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "No Rows",
			args: args{
				query: "SELECT * from STUDENT where code = '1000'",
			},
			wantErr: true,
		},
		{
			name: "One",
			args: args{
				query: "SELECT * from STUDENT",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gendb.QueryRow[sql_test.Student](db, sql_test.StudentBasicBinder, tt.args.query).
				OnError(func(e error) {
					assert.True(t, tt.wantErr)
					assert.Equal(t, e, sql.ErrNoRows)
				}).
				OnSuccess(func(students sql_test.Student) {
					assert.False(t, tt.wantErr)
				})
		})
	}
}

func TestQuery(t *testing.T) {
	db := sql_test.SetupDB(t)

	type args struct {
		binder func(student *sql_test.Student) []any
		query  string
		args   []any
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantCount int
	}{
		{
			name: "simple",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "SELECT id, code, name, program FROM student",
				args:   nil,
			},
			wantCount: 2,
		},
		{
			name: "simple with args",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "SELECT id, code, name, program FROM student WHERE code = ?",
				args:   []any{sql_test.NewIdCode("0001")},
			},
			wantCount: 1,
		},
		{
			name: "query arg count error",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "SELECT id, code, name, program FROM student WHERE id = ?",
				args:   nil,
			},
			wantErr: true,
		},
		{
			name: "binder column mismatch",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "SELECT id FROM student",
				args:   nil,
			},
			wantErr: true,
		},
		{
			name: "empty sql query",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "",
				args:   nil,
			},
			wantErr: true,
		},
		{
			name: "bad sql query",
			args: args{
				binder: sql_test.StudentBasicBinder,
				query:  "SELEC id, code, name, program FROM foo",
				args:   nil,
			},
			wantErr: true,
		},
		{
			name: "broken binder",
			args: args{
				binder: func(s *sql_test.Student) []any { return []any{nil, nil, nil, nil} },
				query:  "SELECT id, code, name, program FROM student",
				args:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gendb.Query[sql_test.Student](db, tt.args.binder, tt.args.query, tt.args.args...).
				OnError(func(e error) {
					fmt.Println(e)
					assert.True(t, tt.wantErr)
				}).
				OnSuccess(func(students []sql_test.Student) {
					assert.False(t, tt.wantErr)
					assert.Equal(t, tt.wantCount, len(students))
				})
		})
	}
}
