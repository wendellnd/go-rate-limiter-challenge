package cache

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() Cache {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return &Redis{
		client,
	}
}

func (r *Redis) GetIncrement(key string) (int, error) {
	val, err := r.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			r.Client.Set(key, 1, 0)
			return 0, nil
		}
		return 0, err
	}

	return strconv.Atoi(val)
}

func (r *Redis) SetExpire(key string, value int, duration int) error {
	return r.Client.Set(key, value, 0).Err()
}

func (r *Redis) Increment(key string) error {
	return r.Client.Incr(key).Err()
}
