progress:
	go run cmd/progress/main.go

run-sender:
	go run cmd/send/main.go "$(lv)" "$(msg)"
run-sender-sample-info: 
	make run-sender lv="info" msg="Everythings worked well."
run-sender-sample-warning: 
	make run-sender lv="warning" msg="Something went wrong."
run-sender-sample-error: 
	make run-sender lv="error" msg="Run. Run. Or it will explode."
run-sender-kernel-error: 
	make run-sender lv="kern.critical" msg="Kernel will explode."
run-sender-origin:
	go run cmd/send/origin/emit_log_topic.go "kern.critical" "Origin error message"

run-receiver:
	go run cmd/receive/main.go info warning error
run-receiver-w-e:
	go run cmd/receive/main.go warning error > cmd/receive/logs/logs_from_rabbit.log
run-receiver-all:
	go run cmd/receive/main.go "#"
run-receiver-kern:
	go run cmd/receive/main.go "kern.*"
run-receiver-crit:
	go run cmd/receive/main.go "*.critical"
run-receiver-kern-or-crit:
	go run cmd/receive/main.go "kern.*" "*.critical"
run-receiver-origin:
	go run cmd/receive/origin/receive_logs_topic.go "kern.*"

run-rpc-client:
	go run cmd/rpc/client/main.go  
run-rpc-server:
	go run cmd/rpc/server/main.go  

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
