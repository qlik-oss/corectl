FROM golang:1.12 AS builder

ADD . /testconnector
WORKDIR /testconnector

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /testconnector .
CMD ["./main"]
EXPOSE 50051
