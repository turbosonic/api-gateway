shell:
	docker run -it -v $(PWD):/go/src/github.com/turbosonic/api-gateway --workdir /go/src/github.com/turbosonic/api-gateway golang:1.10.2 /bin/sh

run: 
	docker run -v $(PWD):/go/src/github.com/turbosonic/api-gateway -p 8080:8080 --workdir /go/src/github.com/turbosonic/api-gateway golang:1.10.2 go run app.go

test:
	go test -v ./tests

build:
	docker build -t turbosonic/api-gateway .