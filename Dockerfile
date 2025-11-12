FROM golang:1.21.1 as builder
RUN apt-get update && apt-get install -y gcc libc6-dev
WORKDIR /go/src/iot-sms
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/iot-sms .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates curl openssl cron
WORKDIR /root/
COPY --from=builder /go/bin/iot-sms .
EXPOSE 1323
CMD ["./iot-sms"]