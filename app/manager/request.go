package manager

import (
	"gopkg.in/resty.v1"
)

var (
	requestManager *RequestManager
)

func GetRequest() *RequestManager {
	return requestManager
}

type RequestManager struct {
	client *resty.Client
}

func (manager *RequestManager) Setup() (err error) {

	client := resty.New()
	client.SetRetryCount(3)

	requestManager = &RequestManager{client: client}

	return nil
}

func (manager *RequestManager) Close() {

}

func (manager *RequestManager) NewRequest() *resty.Request {
	return manager.client.R()
}
