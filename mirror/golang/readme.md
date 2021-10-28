
    export GOPROXY=https://mirrors.aliyun.com/goproxy/
    export GOPROXY=https://goproxy.io
    export GOPROXY=https://goproxy.cn/
    export GOPROXY=GOPROXY=https://goproxy.cn,https://goproxy.io,direct


```bash
# 私有包
GOPRIVATE=*.corp.example.com,rsc.io/private
GOPRIVATE=*.xx.com
GOINSECURE=*.xx.com

# 设置 GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct


GIT_TERMINAL_PROMPT=1

# 以 http 协议取代 ssh 协议拉取代码
git config --global url."ssh://git@gitlab.xxx.com".insteadOf "http://gitlab.xxx.com"
# 以 gitlab ssh协议拉取代码时需要配置的环境变量
GIT_USER=zeusro
GIT_ACCESS_TOKEN=8zG-DptwbduPpcYhSezY
```
