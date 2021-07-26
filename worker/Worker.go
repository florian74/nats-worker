package worker

import "github.com/nats-io/nats.go"

type Worker interface {
	Connect(url string, queueName string, subject string, workerName string, workerType string)
	Work(msg *nats.Msg)
}
