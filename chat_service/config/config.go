package config

type (
	Config struct {
		Database      Database `yaml:"database"`
		Debug         bool     `yaml:"debug"`
		Host          string   `yaml:"host"`
		Port          string   `yaml:"port"`
		AllowOrigins  string   `yaml:"allow_origins"`
		AllowHeaders  string   `yaml:"allow_headers"`
		SecretKey     string   `yaml:"secret_key"`
		PublicMedia   string   `yaml:"public_media"`
		PrivateMedia  string   `yaml:"private_media"`
		MaxUploadSize int64    `yaml:"max_upload_size"`

		// Based on Days
		AccessTokenLifeSpan int64 `yaml:"access_token_lifespan"`
		// Based on Months
		RefreshTokenLifeSpan int64 `yaml:"refresh_token_lifespan"`
	}

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	}
)
