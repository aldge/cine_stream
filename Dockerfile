FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建应用用户和组
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 从本机复制编译好的 release 目录，保持原有目录结构
# 使用 --chown 在复制时直接设置文件所有者
# 注意：构建前需要先在本机执行 ./build.sh 生成 release 目录
COPY --chown=appuser:appuser release /app/

# 验证二进制文件是否存在并设置执行权限
RUN test -f /app/bin/cine_stream || (echo "Error: /app/bin/cine_stream not found! Please run ./build.sh first." && exit 1) && \
    chmod +x /app/bin/cine_stream

# 设置工作目录
WORKDIR /app/

# 设置非 root 用户
USER appuser

# 暴露端口（根据实际需要调整）
EXPOSE 8080

# 运行应用
CMD ["./bin/cine_stream", "--conf=conf/app.yaml"]
