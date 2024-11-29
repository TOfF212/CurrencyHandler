package redis

import (
	"context"
	"fmt"
	"log"
	"myproject/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

type RedisDataBase struct {
	Client *redis.Client
}

func (r *RedisDataBase) NewClient(ctx context.Context, cfg Config) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
	}
	r.Client = db

}

func (r *RedisDataBase)SetCurrencies(currencies []models.Currency) {
	for _, curr := range currencies {
		if err := r.Client.Set(context.Background(), curr.Currency, curr.Rate, 0).Err(); err != nil {
			log.Printf("failed to set data, error: %s", err.Error())
		}
		fmt.Printf("%s set in redis\n", curr.Currency)
	}
}
