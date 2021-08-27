FROM golang:1.16 as builder
WORKDIR /go/src/github.com/etiennecoutaud/go-demo
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o go-demo

FROM gcr.io/distroless/static-debian10
COPY --from=builder /go/src/github.com/etiennecoutaud/go-demo/go-demo /

EXPOSE 8080
USER 1001
CMD ["./go-demo"]