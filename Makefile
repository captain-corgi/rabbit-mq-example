progress:
	go run cmd/progress/main.go

run-sender:
	go run cmd/send/main.go $(msg)

run-receiver:
	go run cmd/receive/main.go

rabbit-mq-start:
	brew services start rabbitmq

rabbit-mq-stop:
	brew services stop rabbitmq

rabbit-mq-restart:
	brew services restart rabbitmq

rabbit-mq-list:
	rabbitmqctl list_queues

rabbit-mq-list-debug:
	rabbitmqctl list_queues name messages_ready messages_unacknowledged
