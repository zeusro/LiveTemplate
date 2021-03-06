

## 服务管理

```bash
service --status-all


Usage: /etc/init.d/mysql start|stop|restart|reload|force-reload|status

Usage: service < option > | --status-all | [ service_name [ command | --full-restart ] ]

```

[Ubuntu Service系统服务说明与使用方法](http://www.mikewootc.com/wiki/linux/usage/ubuntu_service_usage.html)

## 修改SSH端口


```
sed -i "s/#Port .*/Port 12306/g" /etc/ssh/sshd_config
sed -i "s/Port .*/Port 12306/g" /etc/ssh/sshd_config
firewall-cmd --permanent --zone=public --add-port=12306/tcp
firewall-cmd --reload
service sshd restart

/etc/ssh/sshd_config
```

1. [Changing Your SSH Port For Extra Security on CentOS 6 or 7](https://www.vultr.com/docs/changing-your-ssh-port-for-extra-security-on-centos-6-or-7)
1. [Ubuntu 16.04修改ssh端口](https://www.jianshu.com/p/d88d4c6581f5)

## 升级系统

```bash
sudo apt-get update -y
sudo apt-get upgrade  -y
sudo apt-get dist-upgrade  -y
do-release-upgrade
```

## 升级到 Ubuntu18 之后 ssh服务不可用

执行sshd -t

之后排查配置即可,一般是 ssh 配置过老导致

参考:
[如何将Ubuntu升级到18.04最新版](https://cloud.tencent.com/developer/article/1174343)

