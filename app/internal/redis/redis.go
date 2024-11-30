package redis

import (
	"context"
	"fmt"
	"log"
	"myproject/internal/models"
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
		Addr:        "localhost:6379",
		Password:    "password",
		User:        "user",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
	r.ctx = context.Background()
}

func (r *RedisDataBase) Open() {
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:         r.cfg.Addr,
		Password:     r.cfg.Password,
		DB:           r.cfg.DB,
		Username:     r.cfg.User,
		MaxRetries:   r.cfg.MaxRetries,
		DialTimeout:  r.cfg.DialTimeout,
		ReadTimeout:  r.cfg.Timeout,
		WriteTimeout: r.cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
	}
	r.Client = db
}

func (r *RedisDataBase) Close() {
	r.Client.Close()
}

func (r *RedisDataBase) IsEmpty() bool {
	r.Open()
	defer r.Close()
	keyCount, err := r.Client.DBSize(r.ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка при проверке размера Redis: %v", err)
	}
	if keyCount==0{
		return true
	}
	return false
	
}

func (r *RedisDataBase) SetCurrencies(currencies []models.Currency) {
	r.Open()
	defer r.Close()
	for i, curr := range currencies {
		if err := r.Client.Set(r.ctx, curr.Currency, curr.Rate, 24*time.Second).Err(); err != nil {
			log.Printf("failed to set data, error: %s", err.Error())
		}
		fmt.Printf("%d %s set in redis\n", i, curr.Currency)
	}
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

	}
	return rate, nil
}
