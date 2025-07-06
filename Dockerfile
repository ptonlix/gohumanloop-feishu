FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build/gohumanloop-feishu

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/gohumanloop-feishu main.go


FROM alpine

RUN echo http://mirrors.aliyun.com/alpine/v3.10/main/ > /etc/apk/repositories && \
    echo http://mirrors.aliyun.com/alpine/v3.10/community/ >> /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

RUN mkdir -p log conf data

COPY --from=builder /app/gohumanloop-feishu /app/gohumanloop-feishu

# 声明挂载点
VOLUME ["/app/conf", "/app/data"]

CMD ["./gohumanloop-feishu"]
