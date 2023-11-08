build:
	@go build -o bin/goblog

run: build
	@./bin/goblog