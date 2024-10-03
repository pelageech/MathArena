# syntax=docker/dockerfile:1
FROM golang:1.23.1-alpine AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -C cmd -v -o /usr/local/bin/matharena

FROM gcr.io/distroless/static-debian12
WORKDIR /home/nonroot/
USER nonroot:nonroot
COPY --from=build /usr/local/bin/matharena ./

ENTRYPOINT ["./matharena"]