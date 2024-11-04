
test:
	go test ./pkg/...

build:
	go build -o bin/ ./cmd/...
	
run:
	go run ./cmd/...

clean:
	rm -rf bin/

.PHONY: test build run clean
