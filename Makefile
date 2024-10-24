.PHONY: generate
generate:
	go generate ./...

.PHONY: swagger-build
swagger-build:
	swagger generate spec -o ./swagger.yaml --scan-models

.PHONY: swagger
swagger: swagger-build
	swagger serve .\swagger.yaml

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: math-arena
math-arena:
	go build -o out/math-arena ./cmd/...

.PHONY: goose
goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	source .envrc # Or any other preferred way to specify env vars alike
	goose up
