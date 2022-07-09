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
	"gendb"
	"gendb/internal/sql_test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNullable(t *testing.T) {
	db := sql_test.SetupDB(t)

	// Insert a student with a null IdCode
	require.True(t, sql_test.AddStudentToDB(db, sql_test.NullIdCode, "bobby tables ", "CS").Ok())

	type args struct {
		query     string
		checkNull bool
	}
	tests := []struct {
		name       string
		args       args
		wantCount  int
		wantIsNull bool
	}{
		{
			name: "all",
			args: args{
				query: "SELECT id, code, name, program FROM student",
			},
			wantCount: 3,
		},
		{
			name: "nulls",
			args: args{
				query:     "SELECT id, code, name, program FROM student WHERE code is null",
				checkNull: true,
			},
			wantCount:  1,
			wantIsNull: true,
		},
		{
			name: "not nulls",
			args: args{
				query:     "SELECT id, code, name, program FROM student WHERE code is not null",
				checkNull: true,
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gendb.Query[sql_test.Student](db, sql_test.StudentBasicBinder, tt.args.query).
				OnError(func(e error) {
					assert.Fail(t, e.Error())
				}).
				OnSuccess(func(students []sql_test.Student) {
					assert.Equal(t, tt.wantCount, len(students))
					if tt.args.checkNull {
						assert.Equal(t, tt.wantIsNull, students[0].Code.IsNull())
					}
				})
		})
	}
}
