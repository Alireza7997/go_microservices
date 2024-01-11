package global

import (
	_ "embed"
	"microservice/config"
	"microservice/pkg/database"

	"github.com/doug-martin/goqu/v9"
)

//go:embed version
var Version string

//go:embed name
var Name string

// Config
var CFG *config.Config = nil

// Secret key
var SecretKeyBytes []byte = nil

// Default database
var DB database.RelationalDatabaseFunction = nil

// Goqu instance of the database
var GoquDB *goqu.Database = nil

// SQL connections
var AllSQLCons = map[string]database.RelationalDatabaseFunction{}
