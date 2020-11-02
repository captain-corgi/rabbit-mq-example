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
		name       string     = "logs_topic"
		kind       string     = "topic"
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
		exchange  string          = "logs_topic"
		key       string          = severityFrom(os.Args) // Routing key
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
	if (len(args) < 3) || os.Args[2] == "" {
		s.WriteString("hello")
	} else {
		s.WriteString(strings.Join(args[2:], " "))
	}
	return s.String()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const (
	//Info level
	Info string = "info"
	// Warn level
	Warn string = "warning"
	// Error level
	Error string = "error"
)

func severityFrom(args []string) string {
	if len(args) < 2 {
		return Info
	}

	switch strings.Trim(os.Args[1], " ") {
	case Info, Warn, Error:
		return os.Args[1]
	case "":
		return Info
	default:
		return os.Args[1]
	}
}
