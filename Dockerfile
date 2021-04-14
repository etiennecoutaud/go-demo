FROM golang:1.16
WORKDIR /go/src/github.com/etiennecoutaud/go-demo
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o go-demo

EXPOSE 8080
CMD ["/go-demo"]