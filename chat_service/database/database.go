package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alireza/config"
	"github.com/doug-martin/goqu/v9"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	database *sql.DB
	Database *goqu.Database
)

func InitPostgres(db *config.Database) {
	var err error
	database, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Tehran",
		db.Host, db.Port, db.Username, db.DBName, db.Password))

	if err != nil {
		log.Fatal(err)
	}

	Database = goqu.New("postgres", database)
	migrateLatestChanges()
}

func migrateLatestChanges() {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/",
	}
	n, err := migrate.Exec(database, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}
	if n > 0 {
		fmt.Println("\n==Migrations==")
		fmt.Printf("Applied %d migrations!\n", n)
	}
}
