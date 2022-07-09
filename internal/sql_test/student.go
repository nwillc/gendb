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
