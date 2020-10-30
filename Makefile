progress:
	go run cmd/progress/main.go

run-sender:
	go run cmd/send/main.go $(lv) $(msg)
run-sender-sample-info: 
	make run-sender lv="info" msg="Everythings worked well."
run-sender-sample-warning: 
	make run-sender lv="warning" msg="Something went wrong."
run-sender-sample-error: 
	make run-sender lv="error" msg="Run. Run. Or it will explode."

run-receiver:
	go run cmd/receive/main.go info warning error
run-receiver-w-e:
	go run cmd/receive/main.go warning error > cmd/receive/logs/logs_from_rabbit.log

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
