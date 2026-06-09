# 本地已编译好二进制，直接复制到最小运行时镜像
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY bin/yftools-server /app/server
EXPOSE 8000
CMD ["/app/server"]
