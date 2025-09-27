package database

import (
	"context"
	"errors"
	"github.com/onec-tech/bot/pkg/cache"
	"github.com/onec-tech/bot/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type ContextKey string

var (
	ErrNoRows = errors.New("scanning one: no rows in result set")
	c         = new(Config)
)

const (
	ERROR_KEY_EMPTY                   = "cant empty cache key"
	CONTEXT_CACHE_MODE_KEY ContextKey = "cache"
	CONTEXT_CACHE_KEY      ContextKey = "key"
	CONTEXT_TTL_KEY        ContextKey = "ttl"
)

// Config read only.
type Config struct {
	PostgresMasterAddr string `envconfig:"POSTGRES_MASTER_ADDR"`
	PostgresSlaveAddr  string `envconfig:"POSTGRES_SLAVE_ADDR"`
}

var _ IDatabase = (*Database)(nil)
var _ Transaction = (*Tx)(nil)

func New(cache cache.ICache, log logger.Logger, conf ...*Config) (*Database, error) {
	if err := envconfig.Process("", c); err != nil {
		log.Fatalf("db layer failed to load configuration: %s", err)
	}

	dbl := &Database{cache: cache, log: log}

	if len(conf) == 0 {
		conf = append(conf, c)
	}

	writePool, readPool, err := newConn(conf[0].PostgresMasterAddr, conf[0].PostgresSlaveAddr, cache, log)

	if err != nil {
		return nil, err
	}

	dbl.writePool = writePool
	dbl.readPool = readPool

	return dbl, nil
}

func newConn(pgAddrMaster, pgAddrSlave string, cache cache.ICache, log logger.Logger) (writePool, readPool *pgxpool.Pool, err error) {
	log.Info("Connect to master postgresql:", pgAddrMaster)
	masterConfig, err := pgxpool.ParseConfig(pgAddrMaster)
	if err != nil {
		return nil, nil, err
	}
	masterConfig.ConnConfig.PreferSimpleProtocol = true
	writePool, err = pgxpool.ConnectConfig(context.Background(), masterConfig)
	if err != nil {
		return nil, nil, err
	}
	log.Info("Connected...")

	log.Info("Connect to slaves postgresql")
	slaveConfig, err := pgxpool.ParseConfig(pgAddrSlave)
	if err != nil {
		return nil, nil, err
	}
	slaveConfig.ConnConfig.PreferSimpleProtocol = true
	readPool, err = pgxpool.ConnectConfig(context.Background(), slaveConfig)
	if err != nil {
		return nil, nil, err
	}
	log.Info("Connected...")

	return writePool, readPool, nil
}
