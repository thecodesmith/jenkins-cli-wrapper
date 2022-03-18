SOURCES = $(shell find . -name '*.go')

jenkinsw: main.go $(SOURCES)
	go build -o jenkinsw main.go

fmt:
	go fmt ./...

deps:
	go get -u
	go mod tidy -compat 1.17

clean:
	rm -f jenkinsw
