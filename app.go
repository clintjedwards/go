package main

import (
	"fmt"

	"github.com/clintjedwards/goto/config"
	"github.com/clintjedwards/goto/storage"
	"github.com/clintjedwards/goto/storage/bolt"
	"github.com/clintjedwards/goto/storage/redis"
	"go.uber.org/zap"
)

type app struct {
	config  *config.Config
	storage storage.Engine
}

func newApp() *app {

	config, err := config.FromEnv()
	if err != nil {
		zap.S().Fatalw("could not load config", "error", err)
	}

	storage, err := initStorage(storage.EngineType(config.Database.Engine))
	if err != nil {
		zap.S().Fatalw("could not configure storage", "error", err)
	}

	return &app{
		config:  config,
		storage: storage,
	}
}

// initStorage creates a storage object with the appropriate engine
func initStorage(engineType storage.EngineType) (storage.Engine, error) {

	config, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	switch engineType {
	case storage.BoltEngine:

		boltStorageEngine, err := bolt.Init(config.Database.Bolt)
		if err != nil {
			return nil, err
		}

		return &boltStorageEngine, nil
	case storage.RedisEngine:
		redisStorageEngine, err := redis.Init(config.Database.Redis)
		if err != nil {
			return nil, err
		}
		return &redisStorageEngine, nil
	default:
		return nil, fmt.Errorf("storage backend not implemented: %s", engineType)
	}
}
