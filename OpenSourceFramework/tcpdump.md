
yum install -y openssh-clients  tcpdump

## 基础知识

网卡

ifconfig

## tcpdump


tcpdump -i eth0



进入或离开$host的数据包.

host=$(hostname)

tcpdump host $host  -i eth0

tcpdump -i eth0 > l.log

### 实战

```
tcpdump tcp port 6379 -i eth0

tcpdump udp port 31701

tcpdump  port 6379 -i eth0

# 查看往来通讯
tcpdump ip host 172.31.23.182

# 目的是47.106.99.115的单向包
tcpdump  '(dst  host 47.106.99.115)'

# 来源是172.31.21.187的单向包
tcpdump  '(src host 172.31.21.187)'

屏蔽ssh协议抓包

```

https://www.cnblogs.com/ggjucheng/archive/2012/01/14/2322659.html

https://www.kancloud.cn/digest/wireshark/62470

https://blog.csdn.net/xukai871105/article/details/31008635

