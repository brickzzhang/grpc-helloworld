# grpc-helloworld

## What's it for this project ?

This project is a helloworld demo using grpc framework.

## How to run this project ?

1; generate grpc proto and gateway codes:

- using project binary executables:

```bash
make pb
```

- using executables on local machine:

```bash
make lc-pb
```

2; build binary file:

- using project binary executables:

```bash
make build
```

- using executables on local machine:

```bash
make lc-build
```

3; start grpc server:

```bash
make run
```

## Features used from grpc-ecosystem etc

- unary server interceptor for panic recovery.

- extra grpc tags.

- trace support.

- unary, serer-side stream, client-side stream and bidirectional stream support.

- example of customized interceptor for unary server support.

## what's next ?

1. graceful shutdown
