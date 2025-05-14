package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	values map[string]any
}

func New(path string) (*Config, error) {
	cfg, err := load(path)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func load(filepath string) (*Config, error) {
	slog.Debug("Loading config: " + filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {

	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		cfg := &Config{values: make(map[string]any)}
		return cfg, nil
	}

	cfg := &Config{values: make(map[string]any)}

	for key, value := range raw {
		envKey := "WALLCHEMY_" + strings.ToUpper(key)
		if envVal, exists := os.LookupEnv(envKey); exists {
			value = convertType(value, envVal)
		}
		cfg.values[key] = value
	}

	return cfg, nil
}

func convertType(original any, override string) any {
	switch original.(type) {
	case bool:
		if b, err := strconv.ParseBool(override); err == nil {
			return b
		}
	case int:
		if i, err := strconv.Atoi(override); err == nil {
			return i
		}
	case float64:
		if f, err := strconv.ParseFloat(override, 64); err == nil {
			return f
		}
	}
	return override
}

func (c *Config) Override(key string, value any) {
	c.values[key] = value
}

func (c *Config) Overrides(overrides map[string]any) {
	for key, value := range overrides {
		c.values[key] = value
	}
}

func (c *Config) GetString(key string) string {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func (c *Config) GetInt(key string) int {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func (c *Config) GetBool(key string) bool {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		}
	}
	return false
}

func (c *Config) GetStringWithDefault(key, defaultValue string) string {
	if val := c.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func (c *Config) GetIntWithDefault(key string, defaultValue int) int {
	if val := c.GetInt(key); val != 0 {
		return val
	}
	return defaultValue
}

func (c *Config) GetBoolWithDefault(key string, defaultValue bool) bool {
	if _, exists := c.values[key]; exists {
		return c.GetBool(key)
	}
	return defaultValue
}

func (c *Config) FlagOverride(overrides map[string]any) {
	for k, v := range overrides {
		c.values[k] = v
	}
}
