package manager

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

var (
	queueManager *QueueManager
)

type QueueManager struct {
	connect *nats.Conn
	config  QueueConfig
}

type QueueConfig struct {
	Url           string
	MaxReconnects int
}

func GetQueue() *QueueManager {
	return queueManager
}

func (manager *QueueManager) Setup(config QueueConfig) (err error) {

	opts := []nats.Option{
		nats.MaxReconnects(config.MaxReconnects),
	}

	nc, err := nats.Connect(config.Url, opts...)
	if err != nil {
		return err
	}

	queueManager = &QueueManager{connect: nc, config: config}

	return nil
}

func (manager *QueueManager) Close() {

}

func (manager *QueueManager) Subscribe(sub string, msgHandler func(sub string, data []byte) error) (subscription *nats.Subscription, err error) {

	subscription, err = manager.connect.Subscribe(sub, func(msg *nats.Msg) {
		msgHandler(msg.Subject, msg.Data)
	})

	if err != nil {
		return
	}

	return
}

func (manager *QueueManager) Publish(sub string, msg interface{}) (err error) {

	b, err := json.Marshal(msg)
	if err != nil {
		return
	}

	err = manager.connect.Publish(sub, b)
	if err != nil {
		return
	}

	return nil
}
