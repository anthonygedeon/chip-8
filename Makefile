fmt:
	go fmt ./...

clean:
	go clean

lint:
	go vet ./...