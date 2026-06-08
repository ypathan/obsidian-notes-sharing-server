FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY . /app
RUN go build -o file_server

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/file_server .
EXPOSE 8080
CMD ["./file_server"]
