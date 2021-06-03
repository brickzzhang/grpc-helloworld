FROM golang:latest

COPY ./ /usr/local/services/

WORKDIR /usr/local/services/

RUN go build -o boot.bin internal/main.go

CMD ["./boot.bin", "-f", "configs/config.yaml"]
