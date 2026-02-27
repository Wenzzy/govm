package config

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the govm configuration
type Config struct {
	DefaultVersion string            `toml:"default_version"`
	AutoInstall    bool              `toml:"auto_install"`
	InheritVersion bool              `toml:"inherit_version"` // Search parent dirs for go.mod/go.work
	Aliases        map[string]string `toml:"aliases"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgPath string
)

// DefaultConfig returns a new config with default values
func DefaultConfig() *Config {
	return &Config{
		DefaultVersion: "",
		AutoInstall:    true,
		InheritVersion: false, // Only check current directory by default
		Aliases: map[string]string{
			"stable": "",
			"latest": "",
		},
	}
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	var loadErr error
	cfgOnce.Do(func() {
		paths, err := GetPaths()
		if err != nil {
			loadErr = err
			return
		}
		cfgPath = paths.Config

		cfg = DefaultConfig()

		data, err := os.ReadFile(cfgPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Config doesn't exist yet, use defaults
				loadErr = nil
				return
			}
			loadErr = err
			return
		}

		if err := toml.Unmarshal(data, cfg); err != nil {
			loadErr = err
			return
		}

		// Ensure aliases map is initialized
		if cfg.Aliases == nil {
			cfg.Aliases = make(map[string]string)
		}
	})

	return cfg, loadErr
}

// Save saves the configuration to disk
func Save(c *Config) error {
	if cfgPath == "" {
		paths, err := GetPaths()
		if err != nil {
			return err
		}
		cfgPath = paths.Config
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, data, 0644)
}

// Get returns the current configuration (loads if needed)
func Get() *Config {
	c, _ := Load()
	if c == nil {
		return DefaultConfig()
	}
	return c
}

// Reload forces a reload of the configuration
func Reload() (*Config, error) {
	cfgOnce = sync.Once{}
	cfg = nil
	return Load()
}
