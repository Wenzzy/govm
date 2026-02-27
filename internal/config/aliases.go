package config

import (
	"fmt"
	"strings"
)

// ReservedAliases are aliases that have special meaning
var ReservedAliases = []string{"stable", "latest"}

// SetAlias sets an alias for a version
func SetAlias(name, version string) error {
	cfg := Get()
	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}

	// Normalize version (remove 'go' prefix if present)
	version = NormalizeVersion(version)

	cfg.Aliases[name] = version
	return Save(cfg)
}

// RemoveAlias removes an alias
func RemoveAlias(name string) error {
	cfg := Get()
	if cfg.Aliases == nil {
		return nil
	}

	delete(cfg.Aliases, name)
	return Save(cfg)
}

// GetAlias returns the version for an alias
func GetAlias(name string) (string, bool) {
	cfg := Get()
	if cfg.Aliases == nil {
		return "", false
	}

	version, ok := cfg.Aliases[name]
	return version, ok
}

// ListAliases returns all aliases
func ListAliases() map[string]string {
	cfg := Get()
	if cfg.Aliases == nil {
		return make(map[string]string)
	}
	return cfg.Aliases
}

// ResolveVersion resolves a version string which could be an alias
func ResolveVersion(input string) string {
	// First check if it's an alias
	if version, ok := GetAlias(input); ok && version != "" {
		return version
	}
	// Otherwise return normalized version
	return NormalizeVersion(input)
}

// NormalizeVersion removes common prefixes and normalizes the version string
func NormalizeVersion(version string) string {
	// Remove 'go' prefix if present
	version = strings.TrimPrefix(version, "go")
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")
	// Trim whitespace
	version = strings.TrimSpace(version)
	return version
}

// IsAlias checks if a string is a known alias
func IsAlias(name string) bool {
	cfg := Get()
	if cfg.Aliases == nil {
		return false
	}
	_, ok := cfg.Aliases[name]
	return ok
}

// ValidateAliasName checks if an alias name is valid
func ValidateAliasName(name string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}
	if strings.ContainsAny(name, " \t\n/\\") {
		return fmt.Errorf("alias name cannot contain whitespace or path separators")
	}
	return nil
}
