package main

import (
	"flag"
	"github.com/florian74/nats-worker/worker"
	"github.com/nats-io/nats.go"
)

// run nats-server for queue or nats-server -js to enable stream
// create stream with nats stream add <streamname>
func main() {

	url := flag.String("url", nats.DefaultURL, "url of nats server")
	targetSubject := flag.String("target", "default", "queue name if worker is of type queue, reply name topic if queue is of type Reply")
	topic := flag.String("topic", "default", "topic target name")
	workerName := flag.String("name", "worker", "worker name")
	workerType := flag.String("type", "queue", "worker type is Queue, Reply or Stream")
	flag.Parse()

	w := worker.HelloWorker{}
	w.Connect(*url, *targetSubject, *topic, *workerName, *workerType)
}
