

FROM golang:1.17-alpine3.14 AS builder
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    git \
    build-base \
    openssh

WORKDIR /build
COPY . .
RUN chmod -R 777 ./saved

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .
RUN chmod -R 777 app

# Now build the production container
FROM alpine:3.14
RUN apk update \
    && apk add ca-certificates \
    && adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/ /app/

WORKDIR /app
CMD ["./app"]
