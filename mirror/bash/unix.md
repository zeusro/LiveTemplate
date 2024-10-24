

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


1.设置代理

使用 curl，wget，brew等http应用程序会调用http_proxy和https_proxy这两环境变量进行代理，通过下面方式设置：

```
export http_proxy=http://127.0.0.1:8087
export https_proxy=$http_proxy
```
2.取消设置
```
unset http_proxy https_proxy
```
3.快速切换

可以在 ~/.zshrc 或者 ~/.bash_profile 中添加这样的alias：
```
alias goproxy='export http_proxy=http://127.0.0.1:8087 https_proxy=http://127.0.0.1:8087'
alias disproxy='unset http_proxy https_proxy'
```
这样下次就可以很方便地切换proxy了！