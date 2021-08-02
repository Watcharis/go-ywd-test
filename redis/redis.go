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
	addrRedis := os.Getenv("REDIS_URL")
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     addrRedis,
			Password: "",
			DB:       0,
		}),
	}
}

func (re *Redis) ConnectRedis() *redis.Client {
	client := re.Client
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
