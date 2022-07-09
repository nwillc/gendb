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
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const DBName = "sqlite-database.db"

func SetupDB(t *testing.T) *sql.DB {
	t.Helper()
	_ = os.Remove(DBName)
	return genfuncs.NewResultError(sql.Open("sqlite3", DBName)).
		OnError(func(e error) {
			require.Fail(t, "failed creating database", e.Error())
		}).
		OnSuccess(func(db *sql.DB) {
			t.Cleanup(func() {
				_ = db.Close()
			})
			gendb.Exec(db, `CREATE TABLE student (
				"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
				"code" TEXT,
				"name" TEXT,
				"program" TEXT
			);`).
				OnError(func(e error) {
					require.Fail(t, "failed creating schema", e.Error())
				})
			AddStudentToDB(t, db, NewIdCode("0001"), "fred", "masters")
			AddStudentToDB(t, db, NewIdCode("0002"), "barney", "PHD")
		}).
		OrEmpty()
}
