run-unit-test:
	go test ./...

run-happy-flow:
	go build
	./main.exe parse --filepath="test_data/happy_flow_5.csv"
	./main.exe parse --filepath="test_data/happy_flow_15.csv"
	./main.exe parse --filepath="test_data/happy_flow_30.csv"

run-unhappy-flow:
	make unhappy-flow -i

unhappy-flow:
	go build
	./main.exe parse --filepath="test_data/invalid_character.csv"
	./main.exe parse --filepath="test_data/invalid_consumption_count.csv"
	./main.exe parse --filepath="test_data/invalid_header_position_1.csv"
	./main.exe parse --filepath="test_data/invalid_header_position_2.csv"
	./main.exe parse --filepath="test_data/invalid_interval_value.csv"

start-db:
	docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres -p 5432:5432 -d postgres -c log_statement=all

clean-db:
	docker kill postgres
	docker rm postgres

precommit:
	go mod tidy 
	go mod vendor
	golangci-lint run
