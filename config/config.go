package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config read only.
type Config struct {
	ListenHttpPort int    `envconfig:"PORT" default:"8080"`
	ListenGRPCPort int    `envconfig:"GRPC_PORT" default:"8081"`
	PostgresAddr   string `envconfig:"POSTGRES_ADDR" default:""`
	JWTSecret      string `envconfig:"JWT_SECRET" default:"supersecretkey"`
	JWTTTLMinutes  int    `envconfig:"JWT_TTL_MINUTES" default:"60"`
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
