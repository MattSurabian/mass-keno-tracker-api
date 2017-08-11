package redis_cache

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisCacheService struct {
	CacheNamespace string
}

func (rcs *RedisCacheService) GetString(key string) (string, error) {
	redisConn := getRedisConnection()
	defer redisConn.Close()
	return redis.String(redisConn.Do("GET",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
	))
}

func (rcs *RedisCacheService) SetString(key string, value string, ttlSecs int) {
	redisConn := getRedisConnection()
	defer redisConn.Close()

	redisConn.Do(
		"SET",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
		value,
	)
	if ttlSecs != 0 {
		redisConn.Do(
			"EXPIRE",
			fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
			ttlSecs,
		)
	}
}

func (rcs *RedisCacheService) SetObject(key string, object interface{}, ttlSecs int) {
	serializedObject, err := json.Marshal(object)
	if err == nil {
		rcs.SetString(key, string(serializedObject), ttlSecs)
	} else {
		log.Printf("Error Marshaling object: %s", err)
	}
}

func (rcs *RedisCacheService) GetObject(key string, target interface{}) (err error) {
	serializedObject, err := rcs.GetString(key)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(serializedObject), &target)
	if err != nil || len([]byte(serializedObject)) == 0 {
		log.Printf("Problem unmarshaling cached string: %s", err)
		return
	}

	return
}

func (rcs *RedisCacheService) AppendToSet(key string, member string) {
	redisConn := getRedisConnection()
	defer redisConn.Close()

	redisConn.Do(
		"SADD",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
		member,
	)
}

func (rcs *RedisCacheService) GetSet(key string) (set []string, err error) {
	redisConn := getRedisConnection()
	defer redisConn.Close()

	set, err = redis.Strings(redisConn.Do(
		"SMEMBERS",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
	))

	if err != nil {
		return
	}

	return
}

func (rcs *RedisCacheService) Exists(key string) bool {
	redisConn := getRedisConnection()
	defer redisConn.Close()

	exists, _ := redis.Bool(redisConn.Do(
		"EXISTS",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
	))

	return exists
}

func (rcs *RedisCacheService) Bust(key string) {
	redisConn := getRedisConnection()
	defer redisConn.Close()

	redisConn.Do(
		"DEL",
		fmt.Sprintf("%s:%s", rcs.CacheNamespace, key),
	)
}

func (rcs *RedisCacheService) MarkCacheAsVolatile() {
	rcs.SetString("cache-volatile", "true", 0)
}

func (rcs *RedisCacheService) MarkCacheAsNonVolatile() {
	rcs.Bust("cache-volatile")
}

func (rcs *RedisCacheService) IsCacheVolatile() bool {
	return rcs.Exists("cache-volatile")
}
