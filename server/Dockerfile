# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

# 允许 Go 自动下载 go.mod 所需的更高版本工具链
ENV GOTOOLCHAIN=auto

WORKDIR /app

# 缓存依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- Run Stage ----
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /server .
COPY config.yaml .

EXPOSE 8080

CMD ["./server"]
