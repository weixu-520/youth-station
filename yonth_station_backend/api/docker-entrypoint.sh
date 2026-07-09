#!/bin/sh
# 用 envsubst 替换 YAML 中的 ${VAR} 占位符
# 只替换已导出的环境变量，其他保持不变
envsubst < etc/gateway-api.yaml > /tmp/gateway-api.yaml
mv /tmp/gateway-api.yaml etc/gateway-api.yaml
exec ./gateway-api
