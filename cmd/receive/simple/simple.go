package simple

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

//Receiver is a simple receiver
func Receiver() {
	fmt.Println("Hello Receiver")

	// Connect RabbitMQ server
	conn, err := amqp.Dial(MqConnString)
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	// Create channel to interact with mq server
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	// Prevent dispatch two message to one worker at the same time
	var (
		prefetchCount int  = 1
		prefetchSize  int  = 0
		global        bool = false
	)
	err = ch.Qos(prefetchCount, prefetchSize, global)
	failOnError(err, "Failed to set Qos")

	// Define queue to send/receive messages
	var (
		name       string     = "task_queue"
		durable    bool       = true
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
		autoAck  bool   = false // Set to false, if we want to pass incomplete task to others consumers
		noLocal  bool   = false
	)
	msgs, err := ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
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
