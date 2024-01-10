package config

type (
	Config struct {
		CurrentMicroservice   MicroService
		Gateway               MicroService   `yaml:"gateway"`
		MicroServices         []MicroService `yaml:"microservices"`
		Debug                 bool           `taml:"debug"`
		Domain                string         `yaml:"domain"`
		PWD                   string         `yaml:"pwd"`
		AllowOrigins          string         `yaml:"allow_origins"`
		AllowHeaders          string         `yaml:"allow_headers"`
		MaxConcurrentRequests int64          `yaml:"max_concurrent_requests"`
		SecretKey             string         `yaml:"secret_key"`
	}

	MicroService struct {
		// Databases []database.Database
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	}
)
