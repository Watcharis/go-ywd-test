package redis

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	Client *redis.Client
}

func NewConnectRedis() *Redis {
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func (re *Redis) ConnectRedis() *redis.Client {
	url := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	logrus.Info("pong redis ->", pong)
	if err != nil {
		logrus.Errorln("ping Redis Error ->", err)
	}
	return client
}

// func CloneRedis() *redis.Client {
// 	var c *redis.Client
// 	c = ConnectRedis()
// 	fmt.Println("c ->", c)
// 	return c
// }
