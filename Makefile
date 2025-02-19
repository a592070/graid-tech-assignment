.PHONY: build task1 task1-bonus1 task1-bonus2 task2 task3

build:
	go mod tidy
	go build -o build/task1 task1/basic/main.go
	go build -o build/task1-bonus1 task1/bonus1/main.go
	go build -o build/task1-bonus2 task1/bonus2/main.go

task1: build
	./build/task1

task1-bonus1: build
	./build/task1-bonus1

task1-bonus2: build
	./build/task1-bonus2

task2:
	go run task2/main.go $(filter-out $@, $(MAKECMDGOALS))

task3:
	go run task3/main.go

test:
	echo "test"