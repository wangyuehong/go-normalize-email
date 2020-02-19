test:
	go test -race -v ./...

lint:
	@golangci-lint run ./
	@gofmt -s -l -d ./
