

## 软件加速

### 支持socks5加速的软件

如果你用的是ss代理，在当前终端运行以下命令，那么wget curl 这类网络命令都会经过ss代理

```bash
export ALL_PROXY=socks5://127.0.0.1:1080
```


### apt

```bash
cat << EOF >> /etc/apt/apt.conf
Acquire::http::Proxy "http://proxyAddress:port"
EOF
```

## 把代理服务器地址写入shell配置文件.bashrc或者.zshrc

```bash
proxy(){
    # http_proxy=http://userName:password@proxyAddress:port
    export http_proxy=http://127.0.0.1:7890
    export https_proxy=http://127.0.0.1:7890
    echo 'set proxy'
}

noproxy(){
    unset http_proxy
    unset https_proxy
    echo 'unset proxy'
}
```
