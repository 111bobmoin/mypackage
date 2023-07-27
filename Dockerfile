FROM node:17.1-alpine as ui
#!!!!

#设置工作目录
WORKDIR /app

#添加`/app/node_modules/.bin`到$PATH
ENV PATH /app/node_modules/.bin:$PATH

#安装应用依赖
# add app
COPY ui/ .
#RUN sudo apt install npm
RUN npm install --silent

#创建构建
RUN npm run build --omit=dev

#RUN go env -w GO111MODULE=on
#RUN go env -w GOPROXY=https://goproxy.cn,direct
FROM golang:1.18.9 as builder
#!!!!

#将工作目录设置为golang工作空间
WORKDIR /riotpot

# 复制构建并`embed `。UI文件夹中的Go `文件
COPY --from=ui app/build/ ui/build/
COPY --from=ui app/embed.go ui/

# Copy the dependencies into the image
COPY go.mod ./
COPY go.sum ./

# 下载所有依赖项
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

#将所有内容复制到镜像中

#只复制镜像中的应用文件
COPY api api/
COPY cmd cmd/
COPY internal internal/
COPY pkg pkg/
COPY tools tools/

#复制静态文件
COPY statik/ statik/

ADD Makefile .

RUN make build-all

#从Makefile运行命令来构建所有插件

#并构建项目

#——如果您知道您已经准备好构建版本，请在开发时注释此行——

#免责声明:如果你注释了这行代码，请100%确定二进制文件可以在linux上运行

#RUN apt-get update && apt-get install -y openssh-server \
#    && apt-get install -y openssh-client \
#    && apt-get install -y sshpass
#FROM gcr.io/distroless/base-debian10
FROM  gcr.dockerproxy.com/distroless/base-debian10

ENV GIN_MODE=release

WORKDIR /riotpot

#将依赖复制到图像中
COPY --from=builder /riotpot/bin/ ./

# 复制sshpass可执行文件到最终镜像中
#COPY --from=builder /usr/bin/sshpass /usr/bin/sshpass

# API, UI所需。
EXPOSE 2022

CMD ["./riotpot"]
