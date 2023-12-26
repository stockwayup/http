lint: fmt
	golangci-lint run --enable-all --fix

fmt:
	gofmt -w .

gen:
	go generate ./...

build:
	docker build . -t soulgarden/swup:http-0.0.6 --platform linux/amd64
	docker push soulgarden/swup:http-0.0.6
