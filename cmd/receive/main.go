package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	//MqConnString describe connection string to connect RabbitMQ server.
	MqConnString = "amqp://guest:guest@localhost:5672"

	//ContentTypeTextPlain use as content type when publish a message
	ContentTypeTextPlain = "text/plain"
)

func main() {
	fmt.Println("Hello Receiver")

	// Connect RabbitMQ server
	conn, err := amqp.Dial(MqConnString)
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	// Create channel to interact with mq server
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	// Define queue to send/receive messages
	var (
		name       string     = "hello"
		durable    bool       = false
		autoDelete bool       = false
		exclusive  bool       = false
		noWait     bool       = false
		args       amqp.Table = nil
	)
	q, err := ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	failOnError(err, "Failed to declare a queue")

	// Consume a message from queue
	var (
		queue    string = q.Name
		consumer string = ""
		autoAck  bool   = true
		noLocal  bool   = false
	)
	msgs, err := ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	failOnError(err, "Failed to register a consumer")

	// Receive messages whenever it has pushed to queue using a channel
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit, press Ctrl+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
