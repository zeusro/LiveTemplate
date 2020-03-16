
## 误区

bash 加速其实有个误区，不是你加代理就能代理，而且你调用的软件要支持代理才行。

命令行只是使用环境，一般而言是否走http代理由运行的程序自己决定（采用系统代理 或 程序自己指定代理）

所以根本的解决方案是把网络代理装在路由器上面，让路由器根据IPList加速。

## 不支持代理的软件

1. Ping是ICMP协议，不是TCP/UDP协议，Ping不走，也无法走代理。如果你坚持要能Ping通才行，请考虑常规VPN（PPTP/L2PT/IPSec等）

## 解决方案

[MAC/Ubuntu/centos](unix.md)

[Windows](win.md)

参考链接：
1. [让终端走代理的几种方法](https://blog.fazero.me/2015/09/15/%E8%AE%A9%E7%BB%88%E7%AB%AF%E8%B5%B0%E4%BB%A3%E7%90%86%E7%9A%84%E5%87%A0%E7%A7%8D%E6%96%B9%E6%B3%95/)
1. [windows终端命令行下如何使用代理？](https://github.com/shadowsocks/shadowsocks-windows/issues/1489)
1. [命令行走代理的便捷方式](https://juejin.im/post/5e127308e51d4541360ac518)
1. []()
1. []()