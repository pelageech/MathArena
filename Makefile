.PHONY: all generate swagger vendor

generate:
	go generate ./...

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

vendor:
	go mod tidy && go mod vendor

all: vendor generate swagger

math-arena:
	go build -o out/math-arena ./cmd/...