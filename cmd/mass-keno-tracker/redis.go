package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		IdleTimeout: 120 * time.Second,
		MaxActive:   100, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", MASS_KENO_REDIS_ADDRESS)
			if err != nil {
				log.Fatal("REDIS ERROR | ", err, " | Check MASS_KENO_REDIS_ADDRESS environment variable")
			}

			if MASS_KENO_REDIS_PW != "" {
				_, err := c.Do("AUTH", MASS_KENO_REDIS_PW)
				if err != nil {
					log.Fatal("REDIS ERROR | ", err, " | Check MASS_KENO_REDIS_PW environment variable")
				}
			}

			return c, err
		},
	}

}

var redisPool = newPool()

func GetRedisConnection() redis.Conn {
	return redisPool.Get()
}
