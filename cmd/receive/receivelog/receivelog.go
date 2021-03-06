package receivelog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

const (
	//MqConnString describe connection string to connect RabbitMQ server.
	MqConnString = "amqp://guest:guest@localhost:5672"
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
		exName       string     = "logs_topic"
		exKind       string     = "topic"
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

	// Check if user not input level
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)

		// Bind a queue from channel
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")
	}

	// Consume a message from queue
	var (
		cQueue     string     = q.Name
		cConsumer  string     = ""
		cAutoAck   bool       = false // Set to false, if we want to pass incomplete task to others consumers
		cExclusive bool       = true
		cNoLocal   bool       = false
		cNoWait    bool       = false
		cArgs      amqp.Table = nil
	)
	msgs, err := ch.Consume(cQueue, cConsumer, cAutoAck, cExclusive, cNoLocal, cNoWait, cArgs)
	failOnError(err, "Failed to register a consumer")

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
