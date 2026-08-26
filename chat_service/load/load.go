package load

import (
	"log/slog"

	"github.com/Alireza7997/go_microservices/chat_service/global"
	"github.com/Alireza7997/go_microservices/config"
	"github.com/Alireza7997/go_microservices/pkg/loader"
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
