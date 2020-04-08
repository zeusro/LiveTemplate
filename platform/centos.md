
## 日志

journalctl --since="2018-12-11 10:10:30"

## 查看包依赖深度

yum deplist soft_name

修改默认命令别名

```
## 系统为了防止出错,对一些命令设置了别名,这个时候要手动修改
vi ~/.bashrc
#alias cp='cp -i'
# 注释之后就可以愉快地 cp -rf 了
```

## path

```
# 当前用户
 /etc/profile

# Linux系统
~/.bash_profile
```

## 防火墙


```
firewall-cmd --version  # 查看版本
firewall-cmd --help     # 查看帮助

# 查看设置：
firewall-cmd --state  # 显示状态
firewall-cmd --get-active-zones  # 查看区域信息
firewall-cmd --get-zone-of-interface=eth0  # 查看指定接口所属区域
firewall-cmd --panic-on  # 拒绝所有包
firewall-cmd --panic-off  # 取消拒绝状态
firewall-cmd --query-panic  # 查看是否拒绝

firewall-cmd --reload # 更新防火墙规则
firewall-cmd --complete-reload

# 查看所有打开的端口：
firewall-cmd --zone=dmz --list-ports

# 加入一个端口到区域：
firewall-cmd --zone=dmz --add-port=8080/tcp
# 若要永久生效方法同上

# 打开一个服务，类似于将端口可视化，服务需要在配置文件中添加，/etc/firewalld 目录下有services文件夹，这个不详细说了，详情参考文档
firewall-cmd --zone=work --add-service=smtp

# 移除服务
firewall-cmd --zone=work --remove-service=smtp

# 显示支持的区域列表
firewall-cmd --get-zones

# 设置为家庭区域
firewall-cmd --set-default-zone=home

# 查看当前区域
firewall-cmd --get-active-zones

# 设置当前区域的接口
firewall-cmd --get-zone-of-interface=enp03s

# 显示所有公共区域（public）
firewall-cmd --zone=public --list-all

```

https://wangchujiang.com/linux-command/c/firewall-cmd.html

## 升级内核

### 装内核

[update_centos](update_centos.sh)

## 自定义服务自启


nano /usr/lib/systemd/user/shadowsocks.service

```
[Unit]
Description=shadowsocks
Documentation=https://github.com/shadowsocks/shadowsocks
After=network.target remote-fs.target nss-lookup.target


[Service]
Type=forking
ExecStart=/usr/bin/ssserver -c /etc/shadowsocks.json -d start
ExecReload=/usr/bin/ssserver -c /etc/shadowsocks.json -d restart
ExecStop=sudo ssserver -d stop
PrivateTmp=true

[Install]
WantedBy=multi-user.target


```
然后
```bash
systemctl daemon-reload
systemctl enable /usr/lib/systemd/user/shadowsocks.service
```

其他命令

```
systemctl start shadowsocks.service
systemctl restart shadowsocks.service

systemctl status shadowsocks.service
```

https://blog.csdn.net/lishuoboy/article/details/89925957

## 开机source

修改这个文件

~/.bashrc




## bbr
```
modprobe tcp_bbr
echo "tcp_bbr" >> /etc/modules-load.d/modules.conf
echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
sysctl -p
sysctl net.ipv4.tcp_available_congestion_control
sysctl net.ipv4.tcp_congestion_control
```

    
## 安装IPVS管理工具


```
yum install ipvsadm  -y
ipvsadm -ln

```

1. ipvs 为大型集群提供了更好的可扩展性和性能
1. ipvs 支持比 iptables 更复杂的复制均衡算法（最小负载、最少连接、加权等等）
1. ipvs 支持服务器健康检查和连接重试等功能




VxLAN等的隧道技术封装报文

raw表、mangle表、nat表、filter表




[LVS负载均衡之ipvsadm部署安装（安装篇）](https://blog.51cto.com/blief/1743948)


[ipvs 基本介绍](https://www.qikqiak.com/post/how-to-use-ipvs-in-kubernetes/)



[TCP端口状态说明ESTABLISHED、TIME_WAIT](https://blog.csdn.net/zdwzzu2006/article/details/7713499)




[This example shows vxlan nat traversal, using UDP hole punching.](https://gist.github.com/hkwi/9fc7ebc12790ed10ea55ba38e4f86d0e)


[关于VLAN和VXLAN的理解](https://blog.csdn.net/octopusflying/article/details/77609199)

