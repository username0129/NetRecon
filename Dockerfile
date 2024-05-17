FROM alpine:3.19.1
LABEL authors="Arch"

# 替换 APK 源为阿里云镜像
RUN echo "https://mirrors.aliyun.com/alpine/v3.19/main" > /etc/apk/repositories \
    && echo "https://mirrors.aliyun.com/alpine/v3.19/community" >> /etc/apk/repositories

# 安装 Nginx、CURL 及必需的依赖
RUN apk add --no-cache nginx curl ca-certificates tzdata \
    && adduser -D -g 'www' www \
    && mkdir /www \
    && chown -R www:www /var/lib/nginx \
    && chown -R www:www /www

# 安装 Golang
RUN curl -OL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz \
    && rm go1.22.2.linux-amd64.tar.gz

# 设置 Go 环境变量
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建工作目录
WORKDIR /app

# 添加后端代码
ADD ../backend /app/backend

# 构建后端应用
RUN cd /app/backend && go build -ldflags '-s -w' -o backend

# 添加前端代码
ADD ./frontend/dist /www

# 配置 Nginx
COPY ./nginx/default.conf /etc/nginx/http.d/default.conf

# 清理安装后的文件
RUN apk del curl

# 暴露端口
EXPOSE 80 8081

# 启动脚本
CMD ["sh", "-c", "nginx && /app/backend/backend start "]
