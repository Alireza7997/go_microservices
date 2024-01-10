package global

import (
	_ "embed"
)

//go:embed version
var Version string

//go:embed name
var Name string

// Config
