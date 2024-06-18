build:
	@go build -o bin/ndogo cmd/main.go

run: 
	@./bin/ndogo

test: 
	go test -v ./...