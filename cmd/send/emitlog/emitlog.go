package emitlog

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

//EmitLog send message as broadcaster
func EmitLog() {
	fmt.Println("Hello EmitLog")

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
		name       string     = "logs"
		kind       string     = "fanout"
		durable    bool       = true // Set true to persist queue when server stopped
		autoDelete bool       = false
		internal   bool       = false
		noWait     bool       = false
		args       amqp.Table = nil
	)
	err = ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	var (
		exchange  string          = "logs"
		key       string          = ""
		mandatory bool            = false
		immediate bool            = false
		msg       amqp.Publishing = amqp.Publishing{
			ContentType:  ContentTypeTextPlain,
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent,
		}
	)
	err = ch.Publish(exchange, key, mandatory, immediate, msg)
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
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
