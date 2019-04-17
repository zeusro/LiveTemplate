

## 配置linux SSH


上传ssh公钥到$HOME/.ssh/authorized_keys


```bash
vim /etc/ssh/sshd_config
```

取消或者加上以下三行
```
    RSAAuthentication yes
　　PubkeyAuthentication yes
　　AuthorizedKeysFile .ssh/authorized_keys

```
systemctl restart sshd.service

## 超时相关设置

`ClientAliveInterval`指定了服务器端向客户端请求消息的时间间隔, 默认是0，不发送。
而ClientAliveInterval 300表示5分钟发送一次，然后客户端响应，这样就保持长连接了。

```
ClientAliveInterval 300
```

`ClientAliveCountMax`，默认值3。表示服务器发出请求后客户端没有响应的次数达到一定值，就自动断开，正常情况下，客户端不会不响应。

```
ClientAliveCountMax 3
```

ClientAliveInterval 300
ClientAliveCountMax 3
按上面的配置的话，300*3＝900秒＝15分钟，即15分钟客户端不响应时，ssh连接会自动退出。

service sshd restart


## 使用sshpass配合 iterm2实现多远程登录

```
sshpass -p '<pwd>' ssh -o StrictHostKeychecking=no -p22 root@<host>
```

## 参考链接
1. [ssh几个超时参数](http://www.361way.com/ssh-autologout/4679.html)
1. [SSH原理与运用（一）：远程登录](http://www.ruanyifeng.com/blog/2011/12/ssh_remote_login.html)
1. []()

