package worker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type HelloWorker struct {
	workerType string
	natsConn   *nats.Conn
	name       string
}

func (worker *HelloWorker) Work(msg *nats.Msg) {
	fmt.Printf("%s: received msg %s\n", worker.name, string(msg.Data))

	if msg.Reply != "" {
		worker.natsConn.Publish(msg.Reply, []byte(fmt.Sprintf("Received at %s", time.Now().String())))
	}
}

func (worker *HelloWorker) Connect(url string, queueName string, subject string, workerName string, workerType string) {
	var err error

	worker.name = workerName
	worker.workerType = workerType
	worker.natsConn, err = nats.Connect(url)
	if err != nil {
		panic("could not start nats client " + err.Error())
	}

	switch workerType {
	case "queue":
		startQueueListening(worker.natsConn, subject, queueName, worker.Work)
	case "stream":
		startStreamListening(worker.natsConn, subject, queueName, worker.Work)
	}

}

func startQueueListening(conn *nats.Conn, subject string, queueName string, work func(*nats.Msg)) {
	sub, err := conn.QueueSubscribe(subject, queueName, work)
	defer sub.Unsubscribe()
	defer sub.Drain()
	if err != nil {
		panic("could not connect to queue " + err.Error())
	}
	select {}
}

func startStreamListening(conn *nats.Conn, subject string, queueName string, work func(*nats.Msg)) {

	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		panic("could not connect to stream " + err.Error())
	}

	sub, err := js.Subscribe(subject, work)
	defer sub.Unsubscribe()
	defer sub.Drain()
	if err != nil {
		panic("could not subscribe to stream " + err.Error())
	}

	select {}
}
