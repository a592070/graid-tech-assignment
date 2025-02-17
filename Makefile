.PHONY: build task1 task1-bonus1 task1-bonus2

build:
	go build -o build/task1 task1/basic/main.go
	go build -o build/task1-bonus1 task1/bonus1/main.go
	go build -o build/task1-bonus2 task1/bonus2/main.go

task1: build
	./build/task1

task1-bonus1: build
	./build/task1-bonus1

task1-bonus2: build
	./build/task1-bonus2


test:
	echo "test"