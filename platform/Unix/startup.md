# 开机自启

## Linux启动顺序

```
Linux 系统主要启动步骤:
   1. 读取 MBR 的信息,启动 Boot Manager
             Windows 使用 NTLDR 作为 Boot Manager,如果您的系统中安装多个
             版本的 Windows,您就需要在 NTLDR 中选择您要进入的系统。
             Linux 通常使用功能强大,配置灵活的 GRUB 作为 Boot Manager。
   2. 加载系统内核,启动 init 进程
             init 进程是 Linux 的根进程,所有的系统进程都是它的子进程。
   3. init 进程读取 /etc/inittab 文件中的信息,并进入预设的运行级别,
      按顺序运行该运行级别对应文件夹下的脚本。脚本通常以 start 参数启
      动,并指向一个系统中的程序。
             通常情况下, /etc/rcS.d/ 目录下的启动脚本首先被执行,然后是
             /etc/rcN.d/ 目录。例如您设定的运行级别为 3,那么它对应的启动
             目录为 /etc/rc3.d/ 。
   4. 根据 /etc/rcS.d/ 文件夹中对应的脚本启动 Xwindow 服务器 xorg
             Xwindow 为 Linux 下的图形用户界面系统。
   5. 启动登录管理器,等待用户登录
             Ubuntu 系统默认使用 GDM 作为登录管理器,您在登录管理器界面中
             输入用户名和密码后,便可以登录系统。(您可以在 /etc/rc3.d/
             文件夹中找到一个名为 S13gdm 的链接) 安装sysv-rc-conf
```


## centos

在`/etc/init.d/`目录下添加文件,然后on启动就行

使用chkconfig管理

```
chkconfig shadowsocks off
chkconfig shadowsocks --del

```


http://man.linuxde.net/chkconfig

## centos7



修改 /etc/rc.local

```

chmod +x /opt/script/autostart.sh
su - user -c '/opt/script/autostart.sh'
打开/etc/rc.d/rc.local文件，在末尾增加如下内容

#在centos7中，/etc/rc.d/rc.local的权限被降低了，所以需要执行如下命令赋予其可执行权限
chmod +x /etc/rc.d/rc.local

```

通过Systemctl管理service
 
 
 参考
 
 [CentOS 7添加开机启动服务/脚本](https://blog.csdn.net/wang123459/article/details/79063703)
 [systemctl命令](http://man.linuxde.net/systemctl)
 