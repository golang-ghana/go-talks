package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:ravendb@localhost:5432/db_migratn_db?sslmode=disable",
	)

	if err != nil {
		panic(err)
	}

	m.Up()

}
