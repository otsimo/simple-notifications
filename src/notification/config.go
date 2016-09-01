package notification

import "fmt"

const (
	DefaultGrpcPort   = 18844
	DefaultHealthPort = 8080
)

type DriverConfig struct {
	Type     string                 `json:"type" yaml:"type"`
	Provider string                 `json:"provider" yaml:"provider"`
	Config   map[string]interface{} `json:"config" yaml:"config"`
}

type Config struct {
	GrpcPort        int            `json:"-" yaml:"-"`
	HealthPort      int            `json:"-" yaml:"-"`
	TemplatePath    string         `json:"templatePath" yaml:"templatePath"`
	DefaultLanguage string         `json:"defaultLanguage,omitempty" yaml:"defaultLanguage,omitempty"`
	Drivers         []DriverConfig `json:"drivers" yaml:"drivers"`
}

func (c *Config) GetPortString() string {
	return fmt.Sprintf(":%d", c.GrpcPort)
}
func (c *Config) GetHealthPortString() string {
	return fmt.Sprintf(":%d", c.HealthPort)
}

func NewConfig() *Config {
	return &Config{GrpcPort: DefaultGrpcPort, HealthPort: DefaultHealthPort}
}
