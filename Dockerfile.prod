FROM golang:1.8.3 as builder
WORKDIR /Users/bwebster/apps/go/scp-rest-svr
RUN go get -d -v github.com/gorilla/mux
COPY */*.go  .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /Users/bwebster/apps/go/scp-rest-svr/app .
CMD ["./app"]
