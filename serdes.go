package gendb

import (
	"database/sql"
	"database/sql/driver"
)

type (
	// SerDes combines the interfaces needed to serialize and deserialize custom types to and from database.
	SerDes interface {
		sql.Scanner
		driver.Valuer
	}

	// Nullable indicates type supports the database nullable concept.
	Nullable interface {
		IsNull() bool
	}
)
