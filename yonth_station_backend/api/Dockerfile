# ==========================================
# 阶段 1：编译 Go 后端
# ==========================================
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gateway-api .

# ==========================================
# 阶段 2：运行镜像
# ==========================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata gettext
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/gateway-api .
COPY etc/gateway-api.docker.yaml etc/gateway-api.yaml
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8888
ENTRYPOINT ["./docker-entrypoint.sh"]
