FROM golang:1.17 as builder

WORKDIR /opt
COPY . .

RUN go test -v ./...
RUN go build -o app .

FROM alpine:latest
COPY --from=builder /opt/app .
ENTRYPOINT ["./app"]
