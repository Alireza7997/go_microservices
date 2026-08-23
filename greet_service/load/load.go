package load

import (
	"log/slog"

	"microservice/config"
	"microservice/greet_service/global"
	"microservice/pkg/loader"
)

var cfg = &config.Config{}

func init() {
	if err := loader.ParseYaml("../env.yaml", cfg, true); err != nil {
		slog.Error("failed to parse config", "err", err)
		panic(err)
	}

	if _, ok := cfg.Microservices[global.Name]; ok {
		cfg.CurrentMicroservice = cfg.Microservices[global.Name]
	} else {
		slog.Error("microservice definition not found", "name", global.Name)
		panic("microservice definition not found")
	}

	global.CFG = cfg
}
