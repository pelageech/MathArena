FROM alpine:3.20.3
MAINTAINER Gomicro Dev <dev@gomicro.io>

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.22.1/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
COPY ../.envrc_docker /env_for_goose

WORKDIR /migrations
CMD [ "sh", "-c", "source ../env_for_goose && goose up" ]