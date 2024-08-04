package manager

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/caarlos0/env/v6"
)

var (
	configManager *ConfigManager
)

func GetConfig() *ConfigManager {
	return configManager
}

type ConfigManager struct {
	Env EnvVariable
}

type EnvVariable struct {

	// SERVICES
	Port     string `env:"GO_HTTP_PORT,required"`
	Version  string `env:"VERSION" envDefault:"1.2.12"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	// POSTGRES
	PostgresHost                    string `env:"POSTGRES_HOST,required"`
	PostgresPort                    string `env:"POSTGRES_PORT,required"`
	PostgresUsername                string `env:"POSTGRES_USERNAME,required"`
	PostgresPassword                string `env:"POSTGRES_PASSWORD,required"`
	PostgresDatabaseName            string `env:"POSTGRES_DATABASE_NAME,required"`
	PostgresMinConnSize             int32  `env:"POSTGRES_MIN_CONN_SIZE" envDefault:"0"`
	PostgresMaxConnSize             int32  `env:"POSTGRES_MAX_CONN_SIZE" envDefault:"64"`
	PostgresMaxConnIdleTimeBySecond int    `env:"POSTGRES_CONN_IDLE_TIME_BY_SENCOND" envDefault:"1"`
	PostgresMaxConnLifetimeBySecond int    `env:"POSTGRES_CONN_LIFET_TIME_BY_SENCOND" envDefault:"60"`

	// REDIS
	RedisHost     string `env:"REDIS_HOST,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisPrefix   string `env:"REDIS_PREFIX,required"`

	// ELASTICSEARCH
	ElasticsearchIndexPrefix string `env:"ELASTICSEARCH_INDEX_PREFIX,required"`
	ElasticsearchUrl         string `env:"ELASTICSEARCH_URL,required"`

	// NATS
	NatsUrl                      string `env:"NATS_URL,required"`
	NatsMaxReconnects            int    `env:"NATS_MAX_RECONNECTS" envDefault:"-1"`
	NatsAlertByDisconnectMinutes int    `env:"NATS_ALERT_BY_DISCONNECT_MINUTES" envDefault:"2"`

	// S3
	AWSRegion   string `env:"AWS_REGION,required"`
	AWSS3Bucket string `env:"AWS_S3_BUCKET,required"`
}

func (manager *ConfigManager) Setup() (err error) {

	configManager = &ConfigManager{Env: EnvVariable{}}
	err = env.Parse(&configManager.Env)
	if err != nil {
		fmt.Printf("config parse fail, %+v\n", err)
		return err
	}

	err = configManager.Validate()
	if err != nil {
		fmt.Printf("config validate fail, %+v\n", err)
		return err
	}

	return nil
}

func (manager *ConfigManager) Close() {

}

func (manager ConfigManager) Validate() (err error) {

	// SERVICES
	port, err := strconv.ParseUint(manager.Env.Port, 10, 16)
	if err != nil || port <= 0 || port > uint64(65535) {
		err = errors.New("required environment variable \"GO_HTTP_PORT\" should be 0~65535")
		return
	}

	// POSTGRES
	if manager.Env.PostgresHost == "" {
		err = errors.New("required environment variable \"POSTGRES_HOST\" is not set")
		return
	}

	dbPort, err := strconv.ParseUint(manager.Env.PostgresPort, 10, 16)

	if err != nil || dbPort <= 0 || dbPort > uint64(65535) {
		err = errors.New("required environment variable \"POSTGRES_PORT\" should be 0~65535")
		return
	}

	if manager.Env.PostgresUsername == "" {
		err = errors.New("required environment variable \"POSTGRES_USERNAME\" is not set")
		return
	}

	if manager.Env.PostgresDatabaseName == "" {
		err = errors.New("required environment variable \"POSTGRES_DATABASE_NAME\" is not set")
		return
	}

	// REDIS
	if manager.Env.RedisHost == "" {
		err = errors.New("required environment variable \"REDIS_HOST\" is not set")
		return
	}

	if len(manager.Env.RedisPrefix) == 0 {
		err = errors.New("required environment variable \"REDIS_PREFIX\" is not set")
		return
	}

	// ELASTICSEARCH
	if manager.Env.ElasticsearchIndexPrefix == "" {
		err = errors.New("required environment variable \"ELASTICSEARCH_INDEX_PREFIX\" is not set")
		return
	}

	if manager.Env.ElasticsearchUrl == "" {
		err = errors.New("required environment variable \"ES_URL\" is not set")
		return
	}

	// NATS
	if manager.Env.NatsUrl == "" {
		err = errors.New("required environment variable \"NATS_URL\" is not set")
		return
	}

	// S3
	if manager.Env.AWSRegion == "" {
		err = errors.New("required environment variable \"AWS_REGION\" is not set")
		return
	}

	if manager.Env.AWSS3Bucket == "" {
		err = errors.New("required environment variable \"AWS_S3_BUCKET\" is not set")
		return
	}

	return
}
