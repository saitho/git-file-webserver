package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/creasty/defaults"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/saitho/static-git-file-server/utils"
)

var VERSION = "dev"

var GitUpdateModeCache = "cache"
var GitUpdateModeWebhookGitHub = "webhook_github"

const DownloadLocation = "./git_downloads"

type RepoConfig struct {
	Title   string
	Slug    string
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

func (r *RepoConfig) GetDownloadPath() string {
	return path.Join(DownloadLocation, r.Url)
}

func (r *RepoConfig) GetCacheFilePath() string {
	return path.Join(r.GetDownloadPath() + ".cache")
}

func (r *RepoConfig) GetUpdatedTime() int64 {
	cacheFile, _ := ioutil.ReadFile(r.GetCacheFilePath())
	cacheTime, _ := strconv.Atoi(string(cacheFile))
	return int64(cacheTime)
}

func (r *RepoConfig) GetUpdatedTimeObject() time.Time {
	return time.Unix(r.GetUpdatedTime(), 0)
}

func (r *RepoConfig) IsUpToDate() bool {
	// File does not exist
	_, err := os.Stat(r.GetDownloadPath())
	if os.IsNotExist(err) {
		return false
	}

	if r.Update.Mode != GitUpdateModeCache {
		return true
	}

	// Check if cache is up to date (within cacheTime interval)
	return time.Now().Unix() <= (r.GetUpdatedTime() + int64(r.Update.Cache.Time))
}

type Config struct {
	LogLevel string `yaml:"log_level" default:"warning"`
	Git      struct {
		Repositories []*RepoConfig

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

func validateSlugs(cfg *Config) error {
	var slugs []string
	slugBlacklist := []string{"tag", "branch", "webhook"}
	for _, repository := range cfg.Git.Repositories {
		if utils.Contains(slugBlacklist, strings.ToLower(repository.Slug)) {
			return fmt.Errorf("the slug %s is not allowed as it conflicts with internal routes", repository.Slug)
		}
		if utils.Contains(slugs, strings.ToLower(repository.Slug)) {
			return fmt.Errorf("the slug %s is defined multiple times in configuration", repository.Slug)
		}
		slugs = append(slugs, repository.Slug)
	}
	return nil
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
		log.Warnf("Configuration setting Git.CacheTime is deprecated. Use Git.Update.Cache.Time instead.")
	}
	if len(cfg.Git.Repositories) == 0 {
		if len(cfg.Git.Url) > 0 || len(cfg.Git.WorkDir) > 0 || len(cfg.Git.Update.Mode) > 0 || cfg.Git.Update.Cache.Time > 0 || len(cfg.Git.Update.WebHook.GitHub.Secret) > 0 {
			trimmedUrl := strings.TrimSuffix(cfg.Git.Url, ".git")
			urlParts := strings.Split(trimmedUrl, "/")
			slug := strings.Join(urlParts[len(urlParts)-2:], "/")
			cfg.Git.Repositories = []*RepoConfig{{
				Title:   slug,
				Slug:    slug,
				Url:     cfg.Git.Url,
				WorkDir: cfg.Git.WorkDir,
				Update:  cfg.Git.Update,
			}}
			log.Warnf("Configuration settings in Git other than Repositories are deprecated. Define them inside 'repositories' array.")
		}
	}

	if err := validateSlugs(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
