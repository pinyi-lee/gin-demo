package main

import (
	"fmt"
	"gin-demo/app/manager"
	"log"
	"net/http"
	"time"
)

func main() {
	Setup()
	defer Close()
	RunServer()
}

func Setup() {
	var err error

	// CONFIG
	if err = manager.GetConfig().Setup(); err != nil {
		log.Fatalf("config Setup, error:%v", err)
	}

	// LOGGER
	loggerConfig := manager.LoggerConfig{
		LogLevel: manager.GetConfig().Env.LogLevel,
	}
	if err = manager.GetLogger().Setup(loggerConfig); err != nil {
		log.Fatalf("logger Setup, error:%v", err)
	}

	// REQUEST
	if err = manager.GetRequest().Setup(); err != nil {
		log.Fatalf("request Setup, error:%v", err)
	}

	// POSTGRES
	postgresConfig := manager.PostgresConfig{
		Username:                manager.GetConfig().Env.PostgresUsername,
		Password:                manager.GetConfig().Env.PostgresPassword,
		Host:                    manager.GetConfig().Env.PostgresHost,
		Port:                    manager.GetConfig().Env.PostgresPort,
		DatabaseName:            manager.GetConfig().Env.PostgresDatabaseName,
		MinConnSize:             manager.GetConfig().Env.PostgresMinConnSize,
		MaxConnSize:             manager.GetConfig().Env.PostgresMaxConnSize,
		MaxConnIdleTimeBySecond: time.Duration(manager.GetConfig().Env.PostgresMaxConnIdleTimeBySecond),
		MaxConnLifetimeBySecond: time.Duration(manager.GetConfig().Env.PostgresMaxConnLifetimeBySecond),
	}
	if err = manager.GetPostgres().Setup(postgresConfig); err != nil {
		log.Fatalf("postgres Setup, error:%v", err)
	}

	// CACHE
	cacheConfig := manager.CacheConfig{
		Host:     manager.GetConfig().Env.RedisHost,
		Password: manager.GetConfig().Env.RedisPassword,
		Prefix:   manager.GetConfig().Env.RedisPrefix,
	}
	if err = manager.GetCache().Setup(cacheConfig); err != nil {
		log.Fatalf("cache Setup, error:%v", err)
	}

	// SEARCH
	searchConfig := manager.SearchConfig{
		Url:         manager.GetConfig().Env.ElasticsearchUrl,
		IndexPrefix: manager.GetConfig().Env.ElasticsearchIndexPrefix,
	}
	if err = manager.GetSearch().Setup(searchConfig); err != nil {
		log.Fatalf("search Setup, error:%v", err)
	}

	// NATS
	queueConfig := manager.QueueConfig{
		Url:           manager.GetConfig().Env.NatsUrl,
		MaxReconnects: manager.GetConfig().Env.NatsMaxReconnects,
	}
	if err = manager.GetQueue().Setup(queueConfig); err != nil {
		log.Fatalf("queue Setup, error:%v", err)
	}

	// S3
	s3Config := manager.S3Config{
		AWSS3Region: manager.GetConfig().Env.AWSRegion,
		AWSS3Bucket: manager.GetConfig().Env.AWSS3Bucket,
	}
	if err = manager.GetS3().Setup(s3Config); err != nil {
		log.Fatalf("S3 Setup, error:%v", err)
	}

	// ROUTER
	routerConfig := manager.RouterConfig{
		Version: manager.GetConfig().Env.Version,
	}
	if err = manager.GetRouter().Setup(routerConfig); err != nil {
		log.Fatalf("router Setup, error:%v", err)
	}

	// SCHEDULER
	if err = manager.GetScheduler().Setup(); err != nil {
		log.Fatalf("scheduler Setup, error:%v", err)
	}
	go manager.GetScheduler().Run()

}

func Close() {
	manager.GetConfig().Close()
	manager.GetLogger().Close()
	manager.GetRequest().Close()
	manager.GetPostgres().Close()
	manager.GetCache().Close()
	manager.GetSearch().Close()
	manager.GetQueue().Close()
	manager.GetS3().Close()
	manager.GetRouter().Close()
	manager.GetScheduler().Close()
}

func RunServer() {

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", manager.GetConfig().Env.Port),
		Handler:      manager.GetRouter().GetHandler(),
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe, error:%v", err)
	}
}
