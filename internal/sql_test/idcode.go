package sql_test

import (
	"database/sql/driver"
	"fmt"
	"gendb"
)

var (
	_ gendb.SerDes   = (*IdCode)(nil)
	_ gendb.Nullable = (*IdCode)(nil)
	_ fmt.Stringer   = (*IdCode)(nil)

	// NullIdCode is an empty IdCode.
	NullIdCode = IdCode{valid: false}
)

// IdCode represents a Student ID code. This may be empty.
type IdCode struct {
	code  string
	valid bool
}

func NewIdCode(code string) IdCode {
	return IdCode{code: code, valid: true}
}

func (i IdCode) String() string {
	if i.IsNull() {
		return "NONE"
	}

	return i.code
}

func (i *IdCode) Scan(src any) error {
	if src == nil {
		i.valid = false
		return nil
	}
	switch s := src.(type) {
	case string:
		i.code = s
		i.valid = true
	default:
		return fmt.Errorf("invalid id")
	}
	return nil
}

func (i IdCode) Value() (driver.Value, error) {
	if i.IsNull() {
		return nil, nil
	}
	return i.String(), nil
}

func (i IdCode) IsNull() bool {
	return !i.valid
}
