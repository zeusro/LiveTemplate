FROM alpine:3.8

RUN export version=$(cat /etc/alpine-release) && \
    echo https://mirrors.ustc.edu.cn/alpine/v${version%.*}/main > /etc/apk/repositories && \
    echo https://mirrors.ustc.edu.cn/alpine/v${version%.*}/community >> /etc/apk/repositories;\
    echo "Asia/Shanghai" > /etc/timezone ;\
