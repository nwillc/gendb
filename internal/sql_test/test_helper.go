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

package sql_test

import (
	"database/sql"
	"gendb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/results"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	DBName      = "sqlite-database.db"
	RunExamples = "RUN_EXAMPLES"
)

func SetupDB(t *testing.T) *sql.DB {
	t.Helper()
	return CreateDB().
		OnError(func(e error) {
			require.Fail(t, "failed setting up database", e.Error())
		}).
		OnSuccess(func(db *sql.DB) {
			t.Cleanup(func() {
				_ = db.Close()
			})
		}).OrEmpty()
}

func CreateDB() *genfuncs.Result[*sql.DB] {
	_ = os.Remove(DBName)
	return results.Map[*sql.DB, *sql.DB](
		genfuncs.NewResultError(sql.Open("sqlite3", DBName)),
		func(db *sql.DB) *genfuncs.Result[*sql.DB] {
			return results.Map[sql.Result, *sql.DB](
				gendb.Exec(db, `CREATE TABLE student (
											"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
											"code" TEXT,
											"name" TEXT,
											"program" TEXT
										);`),
				func(_ sql.Result) *genfuncs.Result[*sql.DB] { return genfuncs.NewResult(db) },
			).OnSuccess(func(db *sql.DB) {
				AddStudentToDB(db, NewIdCode("0001"), "fred", "masters")
				AddStudentToDB(db, NewIdCode("0002"), "barney", "PHD")
			})
		},
	)
}

func MaybeRunExamples(t *testing.T) {
	t.Helper()
	if _, ok := os.LookupEnv(RunExamples); !ok {
		t.Skip("skipping, environment variable not set:", RunExamples)
	}
}
