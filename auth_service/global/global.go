package global

import (
	_ "embed"
	"microservice/config"
	"microservice/pkg/database"
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

// SQL connections
var AllSQLCons = map[string]database.RelationalDatabaseFunction{}
