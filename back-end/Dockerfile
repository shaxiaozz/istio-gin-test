FROM golang:1.17.8-alpine as builder
WORKDIR /data/istio-gin-test-code
ENV GOPROXY=https://goproxy.cn
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache upx ca-certificates tzdata
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o istio-gin-test

FROM alpine:3.12 as runner
LABEL maintainer="shaxiaozz <2309108459@qq.com>"
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/istio-gin-test-code/istio-gin-test /istio-gin-test
EXPOSE 8080
CMD ["/istio-gin-test"]
