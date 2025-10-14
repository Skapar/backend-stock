package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config read only.
type Config struct {
	ListenHttpPort int    `envconfig:"PORT" default:"8080"`
	PostgresAddr   string `envconfig:"POSTGRES_ADDR" default:""`
}

// New Config constructor.
func New() *Config {
	return &Config{}
}

// Init initialization from environment variables
func (r *Config) Init() {
	if err := envconfig.Process("", r); err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}
}
