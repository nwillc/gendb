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
