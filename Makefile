build:
	@go build -o bin/app cmd/api/main.go

run: clean build
	@./bin/app

clean:
	@rm -f bin/app

