#!/bin/bash
export PATH=$(go env GOPATH)/bin:$PATH

export GO_HTTP_PORT='9999'
export LOG_LEVEL='INFO'

export POSTGRES_HOST='localhost'
export POSTGRES_PORT='5432'
export POSTGRES_USERNAME='postgres'
export POSTGRES_PASSWORD='postgres'
export POSTGRES_DATABASE_NAME='postgres'

export REDIS_HOST='localhost:6379'
export REDIS_PASSWORD=''
export REDIS_PREFIX='pinyi_'

export ELASTICSEARCH_INDEX_PREFIX='pinyi_'
export ELASTICSEARCH_URL='http://localhost:9200'

export NATS_URL='nats://127.0.0.1:4222'

export AWS_REGION='ap-northeast-1'
export AWS_S3_BUCKET='a bucket name'
