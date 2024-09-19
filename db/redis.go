package db

import "github.com/redis/go-redis/v9"

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		//Addr:     "0.0.0.0:6379", for local run
		Password: "12345",
		DB:       0,
	})
	return client
}
