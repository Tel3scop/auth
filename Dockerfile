FROM golang:1.21-alpine AS builder

COPY . /
WORKDIR /

RUN go mod download
RUN go build -o ./bin/server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/server .

CMD ["./server"]