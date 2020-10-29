package receivelog

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	//MqConnString describe connection string to connect RabbitMQ server.
	MqConnString = "amqp://guest:guest@localhost:5672"

	//ContentTypeTextPlain use as content type when publish a message
	ContentTypeTextPlain = "text/plain"
)

//ReceiveLog is a receiver for logging
func ReceiveLog() {
	fmt.Println("Hello ReceiveLog")

	// Connect RabbitMQ server
	conn, err := amqp.Dial(MqConnString)
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	// Create channel to interact with mq server
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	// Declare exchange for channel
	var (
		exName       string     = "logs"
		exKind       string     = "fanout"
		exDurable    bool       = true // Set true to persist queue when server stopped
		exAutoDelete bool       = false
		exInternal   bool       = false
		exNoWait     bool       = false
		exArgs       amqp.Table = nil
	)
	err = ch.ExchangeDeclare(exName, exKind, exDurable, exAutoDelete, exInternal, exNoWait, exArgs)
	failOnError(err, "Failed to declare an exchange")

	// Define queue to send/receive messages
	var (
		name       string     = ""
		durable    bool       = false
		autoDelete bool       = false
		exclusive  bool       = true
		noWait     bool       = false
		args       amqp.Table = nil
	)
	q, err := ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	failOnError(err, "Failed to declare a queue")

	// Consume a message from queue
	var (
		queue    string = q.Name
		consumer string = ""
		autoAck  bool   = false // Set to false, if we want to pass incomplete task to others consumers
		noLocal  bool   = false
	)
	msgs, err := ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	failOnError(err, "Failed to register a consumer")

	// Bind a queue from channel
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	// Receive messages whenever it has pushed to queue using a channel
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			runTaskByDot(&d)
			log.Printf("Done")
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit, press Ctrl+C")
	<-forever
}

func runTaskByDot(d *amqp.Delivery) {
	// Sleep 1s every dot
	dotCount := bytes.Count(d.Body, []byte("."))
	for i := dotCount; i > 0; i-- {
		log.Printf("This task will finish in %d seconds", i)
		time.Sleep(1 * time.Second)
	}

	// Mark this message already finished, and delete from queue
	d.Ack(false)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
