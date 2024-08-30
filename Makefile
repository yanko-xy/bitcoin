
build:
	@go build -o ./bin/bitcoin

run: build
	@./bin/bitcoin