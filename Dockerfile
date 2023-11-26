# 使用官方 Go 镜像作为构建环境
FROM golang:1.21 as builder

# 设置工作目录
WORKDIR /app

# 将 Go 模块文件复制到容器中
COPY go.mod ./
COPY go.sum ./

# 下载 Go 模块依赖
RUN go mod download

# 将源代码复制到容器中
COPY api/ api/
COPY cmd/ cmd/
COPY pkg/ pkg/

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /main cmd/main/main.go

# 使用 scratch 作为最小运行时环境
FROM scratch

# 从构建器中复制可执行文件到当前目录
COPY --from=builder /main ./

# 运行可执行文件
ENTRYPOINT ["./main"]