
## Windows 注意事项

cmd命令行:(不用socks5)(临时设置)(也可放置环境变量)

```bash
set http_proxy=http://127.0.0.1:1080
set https_proxy=http://127.0.0.1:1080
```

如果exe不是走系统代理的话，就要使用VPN等全局的工具了，这样系统所有的流量都会走代理。

没有VPN的话，使用一个叫SSTap的软件也可以将ss的代理转为网卡级别的全局代理（利用的是OpenVPN提供的虚拟网络适配器）：https://www.sockscap64.com/sstap/ 注意：这个软件是免费的，但并不开源（我没有找到源代码），所以请谨慎使用！

## 参考链接
1. [windows终端翻墙，简易方式](https://gist.github.com/dreamlu/cf7cbc0b8329ac145fa44342d6a1c01d)