install:
	go get -u
lint:
	go fmt
test:
	go test
build:
	go mod tidy
	go mod verify
	go get -u
	go build .
clean:
	rm *.ics
