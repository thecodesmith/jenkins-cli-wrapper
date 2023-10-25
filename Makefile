SOURCES = $(shell find . -name '*.go')

jenkinsw: main.go $(SOURCES)
	go build -o jenkinsw main.go

fmt:
	go fmt ./...

deps:
	go get -u
	go mod tidy -compat 1.17

jenkins-up:
	docker run --name jenkins -d -v jenkins_home:/var/jenkins_home -p 8080:8080 -p 50000:50000 jenkins/jenkins:lts-jdk11

jenkins-down:
	docker stop jenkins

setup:
	go install github.com/spf13/cobra-cli@v1.3.0

clean:
	rm -f jenkinsw
