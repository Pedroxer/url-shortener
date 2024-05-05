FROM golang:1.21-alpine as builder
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go build -o  app ./cmd/main.go
CMD ["./app"]



