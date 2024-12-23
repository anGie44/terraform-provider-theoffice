default: fmt lint install generate

build:
	go build -v ./...

install:
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout 120s -parallel=10 ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

PHONY: fmt lint test testacc build install generate
