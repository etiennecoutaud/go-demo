FROM golang:1.16
WORKDIR /go/src/github.com/etiennecoutaud/go-demo
COPY . .

RUN go build -ldflags="-s" -o go-demo
FROM scratch
COPY --from=0 /go/src/github.com/etiennecoutaud/go-demo/go-demo /
EXPOSE 8080
CMD ["/go-demo"]