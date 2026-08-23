package global

import (
	_ "embed"
	"microservice/config"
)

//go:embed version
var Version string

//go:embed name
var Name string

// Config
var CFG *config.Config = nil
