package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	once   sync.Once
)

type RedisConfig struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	DB       int    `json:"db,omitempty" yaml:"db,omitempty"`
}

func NewDefaultRedisCfg() RedisConfig {
	return RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "123456",
		DB:       0,
	}
}

func MustInit(redisConf RedisConfig) {
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
		// todo check redis config
		Client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: redisConf.Password,
			DB:       redisConf.DB,
		})
		_, err := Client.Ping(context.Background()).Result()
		if err != nil {
			panic(err.Error())
		}
	})
}
