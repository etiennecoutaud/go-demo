FROM golang:1.10
WORKDIR /go/src/github.com/etiennecoutaud/go-demo
COPY ..
RUN go get github.com/gorilla/mux
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o go-demo
FROM scratch
COPY --from=0 /go/src/github.com/etiennecoutaud/go-demo/go-demo /
EXPOSE 8080
CMD ["/go-demo"]
