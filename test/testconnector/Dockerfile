FROM golang:1.9 AS builder
RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/qlik-oss/corectl/test/testconnector
WORKDIR /go/src/github.com/qlik-oss/corectl/test/testconnector

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/qlik-oss/corectl/test/testconnector .
CMD ["./main"]
EXPOSE 50051
