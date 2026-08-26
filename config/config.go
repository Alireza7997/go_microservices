package config

import "github.com/Alireza7997/go_microservices/pkg/database"

type (
	Config struct {
		CurrentMicroservice   Microservice
		Gateway               Microservice            `yaml:"gateway"`
		Microservices         map[string]Microservice `yaml:"microservices"`
		Debug                 bool                    `yaml:"debug"`
		Domain                string                  `yaml:"domain"`
		PWD                   string                  `yaml:"pwd"`
		AllowOrigins          string                  `yaml:"allow_origins"`
		AllowHeaders          string                  `yaml:"allow_headers"`
		MaxConcurrentRequests int64                   `yaml:"max_concurrent_requests"`
		SecretKey             string                  `yaml:"secret_key"`
	}

	Microservice struct {
		Databases map[string]database.Database
		IP        string `yaml:"ip"`
		Port      string `yaml:"port"`
	}
)
