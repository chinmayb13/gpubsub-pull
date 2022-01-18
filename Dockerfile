# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder

WORKDIR /app
ENV APP_DIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /app
RUN chmod -R 755 /app
CMD ["/app/main"]