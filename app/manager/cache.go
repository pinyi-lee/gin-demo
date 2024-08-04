package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	cacheManager *CacheManager
)

func GetCache() *CacheManager {
	return cacheManager
}

type CacheManager struct {
	client  *redis.Client
	config  CacheConfig
	context context.Context
}

type CacheConfig struct {
	Host     string
	Password string
	Prefix   string
}

func (manager *CacheManager) Setup(config CacheConfig) (err error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       1,
	})

	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		fmt.Printf("ping redis client err\n")
		return err
	}

	cacheManager = &CacheManager{client: client, config: config, context: context.Background()}

	return nil

}

func (manager *CacheManager) Close() {

}

func (manager *CacheManager) Get(key string) (data string, err error) {

	cmd := manager.client.Get(manager.context, manager.config.Prefix+key)
	if cmd == nil {
		err = errors.New("redis client get err")
		return
	}

	data, err = cmd.Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	return
}

func (manager *CacheManager) Set(key string, value interface{}, expiration time.Duration) error {
	return manager.client.Set(manager.context, manager.config.Prefix+key, value, expiration).Err()
}

func (manager *CacheManager) Del(keys ...string) error {

	newKeys := []string{}
	for _, key := range keys {
		newKeys = append(newKeys, manager.config.Prefix+key)
	}

	return manager.client.Del(manager.context, newKeys...).Err()
}

func (manager *CacheManager) Scan(key string, size int64) (keyList []string, err error) {

	var cursor uint64
	for {
		var keys []string
		keys, cursor, err = manager.client.Scan(manager.context, cursor, manager.config.Prefix+key+"*", size).Result()
		if err != nil {
			return []string{}, err
		}
		keyList = append(keyList, keys...)
		if cursor == 0 {
			break
		}
	}

	return
}

func (manager *CacheManager) UnlinkWithoutPrefix(keys ...string) error {

	newKeys := []string{}
	for _, key := range keys {
		newKeys = append(newKeys, manager.config.Prefix+key)
	}

	return manager.client.Unlink(manager.context, newKeys...).Err()
}
