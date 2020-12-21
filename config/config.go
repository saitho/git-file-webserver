package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Git struct {
		Url       string
		WorkDir   string `yaml:"work_dir"`
		CacheTime int    `yaml:"cache_time"`
	}
	Files []string
}

func LoadConfig(configPath string) (Config, error) {
	var cfg Config
	f, err := os.Open(configPath)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
