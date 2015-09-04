package notification

import "fmt"

type DriverConfig struct {
	Type     string                 `json:"type" yaml:"type"`
	Provider string                 `json:"provider" yaml:"type"`
	Config   map[string]interface{} `json:"config" yaml:"type"`
}

type Config struct {
	Port         int            `json:"port" yaml:"type"`
	TemplatePath string         `json:"templatePath" yaml:"type"`
	Drivers      []DriverConfig `json:"drivers" yaml:"type"`
}

func (c *Config) GetPortString() string {
	return fmt.Sprintf(":%d", c.Port)
}

func NewConfig() *Config {
	return &Config{Port: 50051}
}
