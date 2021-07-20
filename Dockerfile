# stage build

FROM golang:alpine
COPY ./ /usr/local/services/
WORKDIR /usr/local/services/
RUN go build -tags netgo -o boot.bin internal/main.go


# stage run

FROM alpine:latest
COPY --from=0 /usr/local/services/boot.bin /usr/local/services/
COPY --from=0 /usr/local/services/configs /usr/local/services/
# copy generated swagger files
COPY --from=0 /usr/local/services/api/swagger /usr/local/services/api/swagger/
COPY --from=0 /usr/local/services/workshop/swagger/index.html /usr/local/services/workshop/swagger/index.html
WORKDIR /usr/local/services/
CMD ["./boot.bin", "-f", "config.yaml"]
