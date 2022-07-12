<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# gendb
Go SQL database utilities leveraging generics.

## Dependencies
This uses [genfuncs](https://github.com/nwillc/genfuncs) for some Go generics and functional code.

# gendb

```go
import "gendb"
```

## Index

- [func Exec(db *sql.DB, query string, args ...any) *genfuncs.Result[sql.Result]](<#func-exec>)
- [func ExecContext(db *sql.DB, ctx context.Context, query string, args ...any) *genfuncs.Result[sql.Result]](<#func-execcontext>)
- [func Query[T any](db *sql.DB, binder func(*T) []any, query string, args ...any) *genfuncs.Result[[]T]](<#func-query>)
- [func QueryContext[T any](db *sql.DB, ctx context.Context, binder func(*T) []any, query string, args ...any) *genfuncs.Result[[]T]](<#func-querycontext>)
- [func QueryRow[T any](db *sql.DB, binder func(*T) []any, query string, args ...any) *genfuncs.Result[T]](<#func-queryrow>)
- [func QueryRowContext[T any](db *sql.DB, ctx context.Context, binder func(*T) []any, query string, args ...any) *genfuncs.Result[T]](<#func-queryrowcontext>)
- [func SingleBinder[T any](t *T) []any](<#func-singlebinder>)
- [type Nullable](<#type-nullable>)
- [type SerDes](<#type-serdes>)


## func [Exec](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L29-L33>)

```go
func Exec(db *sql.DB, query string, args ...any) *genfuncs.Result[sql.Result]
```

Exec calls ExecContext with the default background context.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"fmt"
	"gendb"
	"gendb/internal/sql_test"
)

func main() {
	db := sql_test.CreateDB()
	exec := gendb.Exec(db.OrEmpty(), "UPDATE student set program = ? WHERE code = '0001'", "CS")
	count, _ := exec.OrEmpty().RowsAffected()
	fmt.Println("Updated:", count)
}
```

#### Output

```
Updated: 1
```

</p>
</details>

## func [ExecContext](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L38-L43>)

```go
func ExecContext(db *sql.DB, ctx context.Context, query string, args ...any) *genfuncs.Result[sql.Result]
```

ExecContext executes a query with arguments and returns a sql.Result summarizing the effect of the statement.

## func [Query](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L55-L60>)

```go
func Query[T any](db *sql.DB, binder func(*T) []any, query string, args ...any) *genfuncs.Result[[]T]
```

Query calls QueryContext with the default background context.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"fmt"
	"gendb"
	"gendb/internal/sql_test"
)

func main() {
	db := sql_test.CreateDB()
	type studentProgram struct {
		name    string
		program string
	}
	binder := func(s *studentProgram) []any { return []any{&s.name, &s.program} }
	results := gendb.Query[studentProgram](
		db.OrEmpty(),
		binder,
		"SELECT name, program FROM student")
	for _, student := range results.OrEmpty() {
		fmt.Println(student.name, student.program)
	}
}
```

#### Output

```
fred masters
barney PHD
```

</p>
</details>

## func [QueryContext](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L65-L71>)

```go
func QueryContext[T any](db *sql.DB, ctx context.Context, binder func(*T) []any, query string, args ...any) *genfuncs.Result[[]T]
```

QueryContext  performs a query with arguments using a binder to assign the rows returned to a slice of results.

## func [QueryRow](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L93-L98>)

```go
func QueryRow[T any](db *sql.DB, binder func(*T) []any, query string, args ...any) *genfuncs.Result[T]
```

QueryRow calls QueryRowContext with the default Background context.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"fmt"
	"gendb"
	"gendb/internal/sql_test"
)

func main() {
	db := sql_test.CreateDB()
	results := gendb.QueryRow[int](
		db.OrEmpty(),
		gendb.SingleBinder[int],
		"SELECT count(*) FROM student")
	fmt.Println("Count:", results.OrEmpty())
}
```

#### Output

```
Count: 2
```

</p>
</details>

## func [QueryRowContext](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L104-L110>)

```go
func QueryRowContext[T any](db *sql.DB, ctx context.Context, binder func(*T) []any, query string, args ...any) *genfuncs.Result[T]
```

QueryRowContext performs a query with arguments using a binder to assign the first row returned to a single result. If multiple rows are returned the first one is used. If no rows are returned the error sql.ErrNoRows is returned.

## func [SingleBinder](<https://github.com/nwillc/genfuncs/blob/master/gendb.go#L131>)

```go
func SingleBinder[T any](t *T) []any
```

SingleBinder is a binder for base types.

## type [Nullable](<https://github.com/nwillc/genfuncs/blob/master/serdes.go#L33-L35>)

Nullable indicates type supports the database nullable concept.

```go
type Nullable interface {
    IsNull() bool
}
```

## type [SerDes](<https://github.com/nwillc/genfuncs/blob/master/serdes.go#L27-L30>)

SerDes combines the interfaces needed to serialize and deserialize custom types to and from database.

```go
type SerDes interface {
    sql.Scanner
    driver.Valuer
}
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
