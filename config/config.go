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

type RepoConfig struct {
	Url     string
	WorkDir string `yaml:"work_dir"`
	Update  struct {
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

type Config struct {
	LogLevel string `yaml:"log_level" default:"warning"`
	Git      struct {
		Repositories []RepoConfig

		Url       string // deprecated: Use setting from Repositories instead
		WorkDir   string `yaml:"work_dir"`   // deprecated: Use setting from Repositories instead
		CacheTime int    `yaml:"cache_time"` // deprecated: use Repositories[].Git.Update.Cache.Time instead
		Update    struct {
			Mode  string `default:"cache"` // deprecated: Use setting from Repositories instead
			Cache struct {
				Time int `yaml:"time"` // deprecated: Use setting from Repositories instead
			} // deprecated: Use setting from Repositories instead
			WebHook struct {
				GitHub struct { // deprecated: Use setting from Repositories instead
					Secret string // deprecated: Use setting from Repositories instead
				}
			} `yaml:"webhook"` // deprecated: Use setting from Repositories instead
		} // deprecated: Use setting from Repositories instead
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
	if len(cfg.Git.Repositories) == 0 {
		if len(cfg.Git.Url) > 0 || len(cfg.Git.WorkDir) > 0 || len(cfg.Git.Update.Mode) > 0 || cfg.Git.Update.Cache.Time > 0 || len(cfg.Git.Update.WebHook.GitHub.Secret) > 0 {
			cfg.Git.Repositories = []RepoConfig{{
				Url:     cfg.Git.Url,
				WorkDir: cfg.Git.WorkDir,
				Update:  cfg.Git.Update,
			}}
			fmt.Printf("Configuration settings in Git other than Repositories are deprecated. Define them inside 'repositories' array.")
		}
	}

	return cfg, nil
}
