FROM golang:1.17 as builder

WORKDIR /opt
COPY . .

RUN go test -v ./...
RUN go build -o app .

FROM ubuntu:focal
WORKDIR /opt
COPY --from=builder /opt/app .
RUN mkdir files
CMD ["./app"]


