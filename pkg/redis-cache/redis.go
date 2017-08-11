package redis_cache

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"time"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		IdleTimeout: 120 * time.Second,
		MaxActive:   100, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("REDIS_CACHE_HOST"))
			if err != nil {
				log.Fatal("REDIS ERROR | ", err, " | Check REDIS_CACHE_HOST environment variable")
			}

			redisPass := os.Getenv("REDIS_CACHE_PASS")
			if redisPass != "" {
				_, err := c.Do("AUTH", redisPass)
				if err != nil {
					log.Fatal("REDIS ERROR | ", err, " | Check REDIS_CACHE_PASS environment variable")
				}
			}

			return c, err
		},
	}

}

var redisPool = newPool()

func getRedisConnection() redis.Conn {
	return redisPool.Get()
}
