# CrazySquirrel: Squirrel wrappers for Cookoo

This library provides hooks for using the Squirrel library from within
Cookoo.

[Cookoo](https://github.com/Masterminds/cookoo) is a chain-of-command
application building framework written in Go.

[Squirrel](https://github.com/lann/squirrel) is a database query
building library that is both powerful and easy to use.

## Installation

```
$ go get github.com/technosophos/crazysquirrel/db
```

Or, for [glide](https://github.com/Masterminds/glide):

```yaml
package: your/package/name
import:
  - package: github.com/Masterminds/cookoo
  - package: github.com/lann/builder
  - package: github.com/lann/squirrel
  - package: github.com/technosophos/crazysquirrel
  # ...

```

## Usage

The main purpose of this library is to make it trivially easy to use the
following workflow:

1. Set up a database connection early in your program.
2. Access the database from anywhere that the cookoo.Context is
   available.
3. Provide transparent and simple access to the Squirrel builder and
   Proxy/Runner.

Below is an extended example.

```go
	package main

	import (
		"github.com/Masterminds/cookoo"
		"github.com/technosophos/crazysquirrel/db"
		"database/sql"
	)

	func main() {
		reg, router, cxt := cookoo.Cookoo()

		conn, err := sql.Open("postgres", "something")
		// ...

		// Set up a Postgres-friendly Squirrel environment.
		db.SetupDatasource(cxt, conn, db.Postgres)
	}

	func MyCommand(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
		// Now we can use the database from anywhere that has a Context.

		rows, err := db.Builder(c).Select("*").From("MyTable").Query()

		// ...
	}

	func Example(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
		// Get the *sql.DB
		conn := db.Db(c)

		// Get a squirrel.DbProxy (which happens to also be a squirrel.BaseRunner)
		// This implements most of the operations that are on *sql.Db, but
		// it wraps them in the statement cache.
		cache := db.Runner(c)

		// Get a squirrel.Builder
		builder := db.Builder(c)
	}
```
