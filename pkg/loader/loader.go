package loader

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ParseYaml(path string, cfg interface{}, errorOnFileNotFound bool) error {
	switch filepath.Ext(path) {

	case ".yaml", ".yml":
		file, err := os.Open(path)
		if err != nil {
			if errorOnFileNotFound {
				return err
			} else {
				return nil
			}
		}

		defer file.Close()

		if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
			return err
		}

		return err

	default:
		return errors.New("unknown file extension")
	}
}

func ParseYamlBytes(data []byte, cfg interface{}) error {
	return yaml.Unmarshal(data, cfg)
}

func ReadConfigFile(cfg interface{}) error {
	if err := ParseYaml("env.yaml", cfg, false); err != nil {
		return err

	} else if err := ParseYaml("env.yml", cfg, false); err != nil {
		return err
	}

	return nil
}
