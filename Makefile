build:
	go build -o bin/notifier


run: build
	./bin/notifier

test:
	go test -v ./...