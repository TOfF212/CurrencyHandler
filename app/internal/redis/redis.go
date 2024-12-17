package redis

import (
	"api/internal/models"
	"context"
	"fmt"
	"log"

	"strconv"
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
	ctx    context.Context
	cfg    Config
}

func (r *RedisDataBase) Init() {
	r.cfg = Config{
		Addr:     "redis_db:6379",
		Password: "",
		DB:       0,
	}
	r.ctx = context.Background()
}

func (r *RedisDataBase) Open() {
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:     r.cfg.Addr,
		Password: r.cfg.Password,
		DB:       r.cfg.DB,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
	}
	r.Client = db
}

func (r *RedisDataBase) Close() {
	r.Client.Close()
}

func (r *RedisDataBase) SetCurrency(curr models.Currency) {
	r.Open()
	defer r.Close()
	if err := r.Client.Set(r.ctx, curr.Currency, curr.Rate, 24*time.Hour).Err(); err != nil {
		log.Printf("failed to set data, error: %s", err.Error())
	}
	fmt.Printf("%s set in redis\n", curr.Currency)
}

func (r *RedisDataBase) GateRate(currency string) (float64, error) {
	r.Open()
	defer r.Close()
	val, err := r.Client.Get(r.ctx, currency).Result()
	if err == redis.Nil {
		return 0, models.ErrorCurrencyNotFound
	} else if err != nil {
		log.Print("error get currency from redis")
		return 0, err
	}

	rate, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
