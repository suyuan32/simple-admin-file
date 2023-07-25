FROM golang:1.20.6-alpine3.17 as builder

# Define the project name | 定义项目名称
ARG PROJECT=file

WORKDIR /build
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w" -o /build/${PROJECT}_api ${PROJECT}.go

FROM nginx:1.25.0-alpine

# Define the project name | 定义项目名称
ARG PROJECT=file
# Define the config file name | 定义配置文件名
ARG CONFIG_FILE=file.yaml
# Define the author | 定义作者
ARG AUTHOR="yuansu.china.work@gmail.com"

LABEL org.opencontainers.image.authors=${AUTHOR}

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}

COPY --from=builder /build/${PROJECT}_api ./
COPY --from=builder /build/etc/${CONFIG_FILE} ./etc/
COPY deploy/nginx/default.conf /etc/nginx/conf.d/
COPY deploy/nginx/entrypoint.sh /docker-entrypoint.d

RUN ["chmod", "+x", "/docker-entrypoint.d/entrypoint.sh"]

EXPOSE 80
EXPOSE 9102

