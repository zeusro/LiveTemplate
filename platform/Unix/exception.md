# 故障处理类



```
# 该命令会输出系统日志的最后10行
dmesg | tail

dmesg -T

```


uptime

## sar

sar -n DEV 1

sar命令在这里可以查看网络设备的吞吐率。在排查性能问题时，可以通过网络设备的吞吐量，判断网络设备是否已经饱和。如示例输出中，eth0网卡设备，吞吐率大概在22 Mbytes/s，既176 Mbits/sec，没有达到1Gbit/sec的硬件上限。


sar -n TCP,ETCP 1

sar命令在这里用于查看TCP连接状态，其中包括：

active/s：每秒本地发起的TCP连接数，既通过connect调用创建的TCP连接；

passive/s：每秒远程发起的TCP连接数，即通过accept调用创建的TCP连接；

retrans/s：每秒TCP重传数量；


## iostat

r/s, w/s, rkB/s, wkB/s：分别表示每秒读写次数和每秒读写数据量（千字节）。读写量过大，可能会引起性能问题。

await：IO操作的平均等待时间，单位是毫秒。这是应用程序在和磁盘交互时，需要消耗的时间，包括IO等待和实际操作的耗时。如果这个数值过大，可能是硬件设备遇到了瓶颈或者出现故障。

avgqu-sz：向设备发出的请求平均数量。如果这个数值大于1，可能是硬件设备已经饱和（部分前端硬件设备支持并行写入）。

%util：设备利用率。这个数值表示设备的繁忙程度，经验值是如果超过60，可能会影响IO性能（可以参照IO操作平均等待时间）。如果到达100%，说明硬件设备已经饱和。


## netstat查看端口状态


```
netstat -n | grep "^tcp" | awk '{print $6}' | sort | uniq -c | sort -n
```



## CPU

```

vmstat 1

```

r：等待在CPU资源的进程数。这个数据比平均负载更加能够体现CPU负载情况，数据中不包含等待IO的进程。如果这个数值大于机器CPU核数，那么机器的CPU资源已经饱和。

free：系统可用内存数（以千字节为单位），如果剩余内存不足，也会导致系统性能问题。下文介绍到的free命令，可以更详细的了解系统内存的使用情况。

si，so：交换区写入和读取的数量。如果这个数据不为0，说明系统已经在使用交换区（swap），机器物理内存已经不足。

us, sy, id, wa, st：这些都代表了CPU时间的消耗，它们分别表示用户时间（user）、系统（内核）时间（sys）、空闲时间（idle）、IO等待时间（wait）和被偷走的时间（stolen，一般被其他虚拟机消耗）。

上述这些CPU时间，可以让我们很快了解CPU是否出于繁忙状态。一般情况下，如果用户时间和系统时间相加非常大，CPU出于忙于执行指令。如果IO等待时间很长，那么系统的瓶颈可能在磁盘IO。


```
# 查看pid CPU占用
pidstat
```


## 查看80端口连接数

netstat -nat|grep -i "80"|wc -l

```
TCP连接状态详解 
LISTEN： 侦听来自远方的TCP端口的连接请求
SYN-SENT： 再发送连接请求后等待匹配的连接请求
SYN-RECEIVED：再收到和发送一个连接请求后等待对方对连接请求的确认
ESTABLISHED： 代表一个打开的连接
FIN-WAIT-1： 等待远程TCP连接中断请求，或先前的连接中断请求的确认
FIN-WAIT-2： 从远程TCP等待连接中断请求
CLOSE-WAIT： 等待从本地用户发来的连接中断请求
CLOSING： 等待远程TCP对连接中断的确认
LAST-ACK： 等待原来的发向远程TCP的连接中断请求的确认
TIME-WAIT： 等待足够的时间以确保远程TCP接收到连接中断请求的确认
CLOSED： 没有任何连接状态
```

[查看linux中的TCP连接数](https://blog.csdn.net/he_jian1/article/details/40787269 )

## 查看磁盘使用情况

```
du -had 1
df -h /
#  查看具体是哪个进程在占用，
fuser -vm /ceshi/ 
```

## 查看cpu核心数

grep 'model name' /proc/cpuinfo | wc -l



## 查看系统平均负载

```
uptime
w
top
```

分别表示系统在过去1分钟、5分钟、15分钟内运行进程队列中的平均进程数量。

没有等待IO，没有WAIT，没有KILL的进程通通都进这个运行队列。


[Linux系统平均负载3个数字的含义](http://www.slyar.com/blog/linux-load-average-three-numbers.html)

## 网络类问题

### 使用dig命令解析域名


```bash
dig @114.114.114.114 baidu.com
dig @114.114.114.114 baidu.com +trace
```

```other
➜  ~ dig @114.114.114.114 baidu.com
; <<>> DiG 9.9.7-P3 <<>> @114.114.114.114 baidu.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 10643
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;baidu.com.			IN	A

;; ANSWER SECTION:
baidu.com.		355	IN	A	123.125.115.110
baidu.com.		355	IN	A	220.181.57.216

;; Query time: 20 msec
;; SERVER: 114.114.114.114#53(114.114.114.114)
;; WHEN: Mon May 07 21:25:10 CST 2018
;; MSG SIZE  rcvd: 70
```

### [ss 命令](https://www.cnblogs.com/peida/archive/2013/03/11/2953420.html)

```bash
# 显示TCP连接
ss -t -a
# UDP
ss -u -a
# 显示 Sockets 摘要
ss -s
# 显示所有状态为Established的HTTP连接
ss -o state established '( dport = :http or sport = :http )' 


```

## Linux临时或永久修改DNS

`sudo vim /etc/resolv.conf`,注意,`nameserver`有数量限制.

```
nameserver 8.8.8.8 #修改成你的主DNS
nameserver 8.8.4.4 #修改成你的备用DNS
search localhost #你的域名
```

## [strace分析系统调用](https://www.cnblogs.com/ggjucheng/archive/2012/01/08/2316692.html)

strace -p 19532

pstack 跟踪进程栈

ipcs 查询进程间通信状态

vmstat 监视内存使用情况

sar 找出系统瓶颈的利器

## [iproute2](https://blog.csdn.net/astrotycoon/article/details/52317288)

取代ifconfig，arp，route，netstat的工具。