package notification
import "fmt"


type DriverConfig struct {
	Type     string "type"
	Provider string "provider"
	Config   map[string]interface{} "config"
}

type Config struct {
	Port         int "port"
	TemplatePath string "templatePath"
	Drivers      []DriverConfig "drivers"
}

func (c *Config)GetPortString() string {
	return fmt.Sprintf(":%d", c.Port)
}

func NewConfig() *Config {
	return &Config{Port:50051}
}