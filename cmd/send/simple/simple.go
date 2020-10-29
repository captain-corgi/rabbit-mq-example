package simple

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

const (
	//MqConnString describe connection string to connect RabbitMQ server.
	MqConnString = "amqp://guest:guest@localhost:5672"

	//ContentTypeTextPlain use as content type when publish a message
	ContentTypeTextPlain = "text/plain"
)

//Sender send message in simple ways
func Sender() {
	fmt.Println("Hello SimpleSender")

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
		name       string     = "task_queue"
		durable    bool       = true // Set true to persist queue when server stopped
		autoDelete bool       = false
		exclusive  bool       = false
		noWait     bool       = false
		args       amqp.Table = nil
	)
	q, err := ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	failOnError(err, "Failed to declare a queue")

	// Publish a message to queue
	body := bodyFrom(os.Args)
	var (
		exchange  string          = ""
		key       string          = q.Name
		mandatory bool            = false
		immediate bool            = false
		msg       amqp.Publishing = amqp.Publishing{
			ContentType:  ContentTypeTextPlain,
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent,
		}
	)
	err = ch.Publish(exchange, key, mandatory, immediate, msg)
	log.Printf("  [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func bodyFrom(args []string) string {
	s := new(strings.Builder)
	if (len(args) < 2) || os.Args[1] == "" {
		s.WriteString("hello")
	} else {
		s.WriteString(strings.Join(args[1:], " "))
	}
	return s.String()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
