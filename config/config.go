package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

var VERSION = "dev"

var GitUpdateModeCache = "cache"
var GitUpdateModeWebhookGitHub = "webhook_github"

type Config struct {
	LogLevel string `yaml:"log_level" default:"warning"`
	Git      struct {
		Url       string
		WorkDir   string `yaml:"work_dir"`
		CacheTime int    `yaml:"cache_time"` // deprecated: use Git.Update.Cache.Time instead
		Update    struct {
			Mode  string `default:"cache"` // cache or webhook_github
			Cache struct {
				Time int `yaml:"time"`
			}
			WebHook struct {
				GitHub struct {
					Secret string
				}
			} `yaml:"webhook"`
		}
	}
	Display struct {
		Branches struct {
			Filter []string
		}
		Tags struct {
			Filter      []string
			Order       string `default:"desc" yaml:"order"`
			ShowDate    bool   `default:"true" yaml:"show_date"`
			VirtualTags struct {
				EnableSemverMajor bool `default:"false" yaml:"enable_semver_major"`
			} `yaml:"virtual_tags"`
		}
		Index struct {
			ShowBranches bool `default:"true" yaml:"show_branches"`
			ShowTags     bool `default:"true" yaml:"show_tags"`
		}
	}
	Files []string
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	_ = defaults.Set(cfg)

	f, err := os.Open(configPath)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return cfg, err
	}

	// Deprecations
	if cfg.Git.CacheTime > 0 {
		cfg.Git.Update.Cache.Time = cfg.Git.CacheTime
		fmt.Printf("Configuration setting Git.CacheTime is deprecated. Use Git.Update.Cache.Time instead.")
	}

	return cfg, nil
}
