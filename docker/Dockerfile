FROM golang:1.24.2 AS builder

WORKDIR /src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /cloud-service .

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate.linux-amd64 /migrate

FROM alpine:3.19 AS final

LABEL maintainer="Bulatnikov Alexandr"

RUN apk add --no-cache postgresql-client

COPY --from=builder /cloud-service /cloud-service
COPY --from=builder /migrate /usr/local/bin/migrate
COPY ./migrations /migrations
COPY ./docker/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh \
    && chmod +x /usr/local/bin/migrate

ENTRYPOINT ["/entrypoint.sh"]
