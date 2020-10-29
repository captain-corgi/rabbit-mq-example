progress:
	go run cmd/progress/main.go

run-sender:
	go run cmd/send/main.go

run-receiver:
	go run cmd/receive/main.go

start-rabbit-mq-server:
	brew services start rabbitmq

stop-rabbit-mq-server:
	brew services stop rabbitmq
