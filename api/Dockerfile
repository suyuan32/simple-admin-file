FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /home
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && cd api \
    && go build -ldflags="-s -w" -o /home/api/file_api file.go

FROM alpine:latest

LABEL MAINTAINER="RyanSU@yuansu.china.work@gmail.com"

WORKDIR /home

COPY --from=builder /home/api/file_api ./
COPY --from=builder /home/api/etc/file.yaml ./

EXPOSE 9102
ENTRYPOINT ./file_api -f file.yaml