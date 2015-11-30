/* Package db provides database helpers for using Squirrel with Cookoo.

This package provides a few convenience tools for consistently working with
Squirrel-based database connections in Cookoo.

Usage:

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

*/
package db

import (
	"database/sql"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/squirrel"
)

const (
	// Key for the raw *sql.Db
	DbObj = "db.Obj"
	// Key for getting the Squirrel StatementCacheProxy.
	DbRunner = "db.Runner"
	// Key for getting the Squirrel StatementBuilder.
	DbBuilder = "db.Builder"
	// Key for getting the DB driver.
	DbDriver = "db.Driver"
)

// Config provides configuration options to SetupDatasource.
//
// Predefined configurations for Postgres and MySQL are provided as variables.
type Config struct {
	Placeholder squirrel.PlaceholderFormat
}

// Predefined config for Postgres defaults
var Postgres = Config{Placeholder: squirrel.Dollar}

// Predifined config for MySQL defaults
var MySQL = Config{Placeholder: squirrel.Question}

// SetupDatasource takes a database connection and creates an easy-to-use
// Squirrel wrapper.
//
// Generally, this is called early in Cookoo's initialization. See the example
// in the package documentation.
//
// The Config struct is used to pass options into the setup. Default
// configurations are provided for Postgres and MySQL.
//
// Example:
//
// 	db.SetupDatasource(c, dbConn, db.MySQL)
func SetupDatasource(c cookoo.Context, db *sql.DB, cfg Config) {

	// Statement Cache is a cache for prepared statements.
	sc := squirrel.NewStmtCacheProxy(db)
	builder := squirrel.StatementBuilder.RunWith(sc)

	if cfg.Placeholder != nil {
		builder.PlaceholderFormat(cfg.Placeholder)
	}

	c.AddDatasource(DbObj, db)
	c.AddDatasource(DbRunner, sc)
	c.AddDatasource(DbBuilder, sc)
}

// Builder fetches the Squirrel statement builder datasource.
//
// Usage:
//
// 	db.Builder(cxt).Select("*").From("foo").Query()
//
func Builder(c cookoo.Context) *squirrel.StatementBuilderType {
	return c.Datasource(DbBuilder).(*squirrel.StatementBuilderType)
}

// Runner fetches a Squirrel DBProxyBeginner datasource.
//
// Usage:
//
// 	db.Runner(cxt).Exec("SELECT * FROM foo LIMIT 10")
func Runner(c cookoo.Context) squirrel.DBProxyBeginner {
	return c.Datasource(DbRunner).(squirrel.DBProxyBeginner)
}

// Db fetches the *sql.DB datasource.
func Db(c cookoo.Context) *sql.DB {
	return c.Datasource(DbObj).(*sql.DB)
}
