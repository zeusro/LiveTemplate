
    export GOPROXY=https://mirrors.aliyun.com/goproxy/
    export GOPROXY=https://goproxy.io
    export GOPROXY=https://goproxy.cn/
    


## go  1.13

```bash
# 私有包
GOPRIVATE=*.corp.example.com,rsc.io/private
go env -w GOPROXY=https://goproxy.cn,direct
```
