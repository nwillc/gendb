package sql_test

import (
	"database/sql"
	"fmt"
	"gendb"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	_ fmt.Stringer = (*Student)(nil)

	StudentBasicBinder = func(s *Student) []any { return []any{&s.Id, &s.Code, &s.Name, &s.Program} }
)

type Student struct {
	Id      int
	Code    IdCode
	Name    string
	Program string
}

func (s Student) String() string {
	return fmt.Sprintf("Student: {id: %d, code: %s, name: %s, program: %s}",
		s.Id, s.Code, s.Name, s.Program)
}

func AddStudentToDB(t *testing.T, db *sql.DB, code IdCode, name string, program string) {
	t.Helper()
	gendb.Exec(
		db,
		"INSERT INTO student(code, name, program) VALUES (?, ?, ?)",
		code, name, program).
		OnError(func(e error) {
			assert.Fail(t, e.Error())
		})
}
