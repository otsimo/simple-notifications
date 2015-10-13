package notification

import "fmt"

const (
	DefaultPort = 18844
)

type DriverConfig struct {
	Type     string                 `json:"type" yaml:"type"`
	Provider string                 `json:"provider" yaml:"provider"`
	Config   map[string]interface{} `json:"config" yaml:"config"`
}

type Config struct {
	Port         int            `json:"port" yaml:"port,omitempty"`
	TemplatePath string         `json:"templatePath" yaml:"templatePath"`
	Drivers      []DriverConfig `json:"drivers" yaml:"drivers"`
}

func (c *Config) GetPortString() string {
	return fmt.Sprintf(":%d", c.Port)
}

func NewConfig() *Config {
	return &Config{Port: DefaultPort}
}
