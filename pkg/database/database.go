package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

type (
	Database struct {
		Type     string `yaml:"type"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		SslMode  string `yaml:"ssl_mode"`
		TimeZone string `yaml:"time_zone"`
		Charset  string `yaml:"charset"`
	}

	RelationalDatabaseFunction func() (*sql.DB, error)
)

func New(dbs map[string]Database, debug bool) (cons map[string]RelationalDatabaseFunction, db RelationalDatabaseFunction, err error) {
	mainOrTest := "main"
	if debug {
		mainOrTest = "test"
	}

	createConnectionFunc := func(dbType, dbConfig string) RelationalDatabaseFunction {
		return func() (*sql.DB, error) {
			db, err := sql.Open(dbType, dbConfig)
			if err != nil {
				return nil, err
			}

			return db, err
		}
	}

	for k, v := range dbs {
		dbConfig := ""

		switch strings.ToLower(v.Type) {
		case "postgres":
			dbConfig = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", v.Host, v.Port, v.Username, v.Password, v.DbName, v.SslMode)
		case "sqlite3":
			if _, err = os.Stat(v.DbName); err != nil {
				_, err = os.Create(v.DbName)
				if err != nil {
					return
				}
			}

			dbConfig = fmt.Sprintf("file:%s?cache=shared&mode=rw", v.DbName)
		case "mysql":
			dbConfig = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", v.Username, v.Password, v.Host, v.Port, v.DbName)
		case "mssql":
			dbConfig = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", v.Host, v.Username, v.Password, v.Port, v.DbName)
		default:
			log.Fatalf("unrecognizable database type : '%s'", v.Type)
		}

		dbFunc := createConnectionFunc(v.Type, dbConfig)

		if mainOrTest == k {
			db = dbFunc
		}

		cons[k] = dbFunc
	}

	return
}

func CloseDbs(cons map[string]*sql.DB) {
	for _, con := range cons {
		con.Close()
	}
}
