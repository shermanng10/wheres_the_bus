.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./where_da_bus

build:
	GOOS=linux GOARCH=amd64 go build -o where_da_bus ./
