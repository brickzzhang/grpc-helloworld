# grpc-helloworld

## What's it for this project ?

This project is a helloworld demo using grpc framework. Only basic grpc unary direction flow is supported now, advanced using of stream grpc will be added later.

## How to run this project ?

1; generate grpc proto and gateway codes:

```bash
make pb
```

2; build binary file:

```bash
go build -o boot.bin internal/main.go
```

3; start grpc server:

```bash
./boot.bin -f configs/config.yaml
```

## Features used from grpc-ecosystem etc

- unary server interceptor for panic recovery.

- extra grpc tags.

- trace support

## what will be next ?

1. grpc stream support
