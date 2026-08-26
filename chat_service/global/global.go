package global

import (
	_ "embed"
	"github.com/Alireza7997/go_microservices/config"
)

//go:embed version
var Version string

//go:embed name
var Name string

// Config
var CFG *config.Config = nil
