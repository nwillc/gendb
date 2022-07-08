package gendb_test

import (
	"gendb"
	"gendb/internal/sql_test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullable(t *testing.T) {
	db := sql_test.SetupDB(t)

	// Insert a student with a null IdCode
	sql_test.AddStudentToDB(t, db, sql_test.NullIdCode, "bobby tables ", "CS")

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
