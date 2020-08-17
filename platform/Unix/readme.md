

## 信号

```
SIGINT:程序终止(interrupt)信号, 在用户键入INTR字符(通常是Ctrl-C)时发出，用于通知前台进程组终止进程。
SIGTERM：进程终止信号，效果等同于*nix shell中不带-9的kill命令；
SIGUSR1：保留给用户使用的信号；
SIGUSR2：同SIGUSR1，保留给用户使用的信号。


```

[Linux信号列表（sigint,sigtstp..)](https://blog.csdn.net/DLUTBruceZhang/article/details/8821690)


 
## 工具类

- xargs

 kgpo | grep partnershop | xargs 'system("echo 1 %"); system("echo 2 %");' 

- date

```
date "+%Y-%m-%d"
date "+%Y-%m-%d %H:%M:%S"

```
https://blog.csdn.net/jk110333/article/details/8590746

* mpv
    ```
    mpv '/Volumes/D/video/[阳光电影www.ygdy8.net].猩球崛起3：终极之战.BD.720p.国英双语.中英双字幕.mkv'  --sub-file='' 
    ```

* 按配置文件启动 redis
    ```
    sudo redis-server /usr/local/Cellar/redis/4.0.6/redis.conf
    ```

- zookeeper
    ```bash
    # To have launchd start zookeeper now and restart at login:
    brew services start zookeeper
    
    # Or, if you don't want/need a background service you can just run:
    zkServer start
    ```

- crontab

    ```
    sudo crontab -e
    0 4 8 * *  bash /home/forzu/rmlog.sh > ~/rm.log
    ```


- ssh
    ```
    #rpm -qa |grep ssh 检查是否装了SSH包
    没有的话yum install openssh-server
    #chkconfig --list sshd 检查SSHD是否在本运行级别下设置为开机启动
    #chkconfig --level 2345 sshd on  如果没设置启动就设置下.
    #service sshd restart  重新启动
    #netstat -antp |grep sshd  看是否启动了22端口.确认下.
    #iptables -nL  看看是否放行了22口.
    #setup---->防火墙设置   如果没放行就设置放行.
    ```
- vi
    ```
    删除一行：dd
    dd 删除一行
    d0 删至行首。
    d$ 删至行尾。
    ^R 恢复u的操作
    查看编码
    :set fileencoding 
    :set fileencoding=utf-8
    ```


### 文件类

- grep


```bash
grep -r "搜索内容" ./
```

- chmod

```
chmod 777 file

r ：读权限，用数字4表示
w ：写权限，用数字2表示
x ：执行权限，用数字1表示
- ：删除权限，用数字0表示
s ：特殊权限 
```
Linux下权限的粒度有 拥有者 、群组 、其它组 三种.

参考:
[Linux权限详解（chmod、600、644、666、700、711、755、777、4755、6755、7755）](https://blog.csdn.net/u013197629/article/details/73608613)



- wc


    wc pom.xml
    行数 单词数 字节数 文件名



* nl

 
    nl package.json

-b  ：指定行号指定的方式，主要有两种：

-b a ：表示不论是否为空行，也同样列出行号(类似 cat -n)；

-b t ：如果有空行，空的那一行不要列出行号(默认值)；

-n  ：列出行号表示的方法，主要有三种：

-n ln ：行号在萤幕的最左方显示；

-n rn ：行号在自己栏位的最右方显示，且不加 0 ；

-n rz ：行号在自己栏位的最右方显示，且加 0 ；

-w  ：行号栏位的占用的位数。

-p 在逻辑定界符处不重新开始计算。 


* tail

```
-f 循环读取

-q 不显示处理信息

-v 显示详细的处理信息

-c<数目> 显示的字节数

-n<行数> 显示行数

--pid=PID 与-f合用,表示在进程ID,PID死掉之后结束. 

-q, --quiet, --silent 从不输出给出文件名的首部 

-s, --sleep-interval=S 与-f合用,表示在每次反复的间隔休眠S秒 
```
    tail -n 5 log2014.log
    tail -f test.log



* find

```
    #查看超过800m文件
    find . -type f -size +800M  -print0 | xargs -0 ls -l
    
　　find /etc -name '*srm*' 
    # 查找在系统中最后10分钟访问的文件　

    #在当前目录查找文件名以一个个小写字母开头，最后是4到9加上.log结束的文件：  
    find . -name "[a-z]*[4-9].log" -print

　  -amin n
　　查找系统中最后N分钟访问的文件
　　-atime n
　　查找系统中最后n*24小时访问的文件
　　-cmin n
　　查找系统中最后N分钟被改变状态的文件
　　-ctime n
　　查找系统中最后n*24小时被改变状态的文件
　　-empty
　　查找系统中空白的文件，或空白的文件目录，或目录中没有子目录的文件夹
　　-false
　　查找系统中总是错误的文件
　　-fstype type
　　查找系统中存在于指定文件系统的文件，例如：ext2 .
　　-gid n
　　查找系统中文件数字组 ID 为 n的文件
　　-group gname
　　查找系统中文件属于gnam文件组，并且指定组和ID的文件
　　
-daystart
　　.测试系统从今天开始24小时以内的文件，用法类似-amin
　　-depth
　　使用深度级别的查找过程方式,在某层指定目录中优先查找文件内容
　　-follow
　　遵循通配符链接方式查找; 另外，也可忽略通配符链接方式查询
　　-help
　　显示命令摘要
　　-maxdepth levels
　　在某个层次的目录中按照递减方法查找
　　-mount
　　不在文件系统目录中查找， 用法类似 -xdev.
　　-noleaf
　　禁止在非UNUX文件系统，MS-DOS系统，CD-ROM文件系统中进行最优化查找
　　-version
　　打印版本数字　　
　　
```
 



* more

```bash
#显示文件中从第3行起的内容
more +3 log2012.log

#从文件中查找第一个出现"day3"字符串的行，并从该处前两行开始显示输出 
more +/day3 log2012.log

#设定每屏显示行数 
more -5 log2012.log
```

参考
[每天一个linux命令(12)：more命令](https://www.cnblogs.com/peida/archive/2012/11/02/2750588.html)


* less


参考
[每天一个linux命令（13）：less 命令](https://www.cnblogs.com/peida/archive/2012/11/05/2754477.html)

### 权限相关
* 文件权限777    
    ```
    cd 你的文件夹路径的上一级目录。
    sudo chmod -R 777 你的文件夹名。
    ```

```
groups <用户名>

```
   
* chown   
```
chown [  -f ] [ -h ] [  -R ] Owner [ :Group ] { File ... | Directory ... }
chown -R  [  -f ] [ -H | -L | -P ] Owner [ :Group ] { File ... | Directory … }
- R 递归式地改变指定目录及其下的所有子目录和文件的拥有者。
- v 显示chown命令所做的工作。
```


### 网络类

#### ip route

```
1）添加到达目标主机的路由记录
   ip route add 目标主机 via 网关
2）添加到达网络的路由记录
   ip route add 目标网络/掩码 via 网关
3）添加默认路由
   ip route add default via 网关
```

https://linux.cn/article-3144-1.html

https://blog.51cto.com/13150617/1963833

https://www.jianshu.com/p/8499b53eb0a5

https://zhuanlan.zhihu.com/p/43279912

http://support.huawei.com/huaweiconnect/enterprise/huawei/m/ViewThread.html?tid=366111

* 查实时流量 iftop

* 查看内网 IP
    ```
    ifconfig -a
    ```
    eth0后面就是内网 IP
    
* 查看外网 IP
    ```
    curl ifconfig.me
    ```


* 统计80端口连接数


    netstat -nat|grep -i "80"|wc -l

* 抓包

    tcpdump

- 统计80端口连接数

    netstat -nat|grep -i "80"|wc -l

* 统计httpd协议连接数


    ps -ef|grep httpd|wc -l


* 转发80端口转发到12306端口

    先修改 pf.conf文件

    ```
    sudo vim /etc/pf.conf
    ```
    按需配置,把rdr on lo0写到rdr-anchor "com.apple/*"下面一行
    ```
    rdr-anchor "com.apple/*"
    rdr on lo0 inet proto tcp from any to 127.0.0.1 port 80 -> 127.0.0.1 port 12306
    ```
    
    ```
   sudo pfctl -d && sudo pfctl -f /etc/pf.conf && sudo pfctl -e

    ```

* 改 host

    terminal 输入
    ```
    sudo vi /etc/hosts 
    # 修改完毕之后先按“esc”，再输入“:wq”

    ```
    
* curl
    ```
    curl -voa https://gz.17zwd.com/
    curl  -voa  $url -H 'Via:xxxx' -H 'Accept-encoding: gzip, deflate, br'
    ```

参考:
[linux curl 命令详解，以及实例](http://blog.51yip.com/linux/1049.html)

* 上传文件命令 [scp](https://www.garron.me/en/articles/scp.html)
```bash
scp [-12346BCEpqrv] [-c cipher] [-F ssh_config] [-i identity_file] [-l limit] [-o ssh_option] [-P port] [-S program] [[user@]host1:]file1 ... [[user@]host2:]file2
```
Examples
```
scp -c blowfish -v /path/to/source-file user@host:/path/to/destination-folder/

```

scp -v a.cap root@47.106.160.166:/root

```
-v 和大多数linux命令中的-v意思一样,用来显示进度。可以用来查看连接、认证、或是配置错误。

-r 递归处理，将指定目录下的文档和子目录一并处理

-C 使能压缩选项

-P 选择端口。注意-p已经被rcp使用

-4 强行使用IPV4地址

-6 强行使用IPV6地址
```

下载文件

```
scp    root@47.106.160.166:/root/a.cap /Volumes/D


```

### 编码类

* enconv 

转换文件编码，比如要将一个GBK编码的文件转换成UTF-8编码，操作如下

    enconv -L zh_CN -x UTF-8 filename

* iconv 

iconv的命令格式如下：
iconv -f encoding -t encoding inputfile
比如将一个UTF-8 编码的文件转换成GBK编码

    iconv -f UTF-8 -t GBK file1 -o file2

### 系统类

    把上一个命令作为cd的目录
    cd !$

### 其他

- cal -y 2013

查看月历

* 进程信息

    ls -l /proc/PID   
    auxv 包含传递给进程的 ELF 解释器信息，格式是每一项都是一个 unsigned long长度的 ID 加上一个 unsigned long 长度的值。
    
    comm包含进程的命令名
    
    cwd符号链接的是进程运行目录；

    exe符号连接就是执行程序的绝对路径；

    cmdline就是程序运行时输入的命令行命令；

    environ记录了进程运行时的环境变量；

    fd目录下是进程打开或使用的文件的符号连接。

参考:

[Linux中 /proc/[pid]目录各文件简析](https://www.hi-linux.com/posts/64295.html)


* 其他好用的工具


    brew - 软件安装
    sox - 命令行录音
    vim － 编辑工具
    tcpdump - 命令行抓包工具
    sed - 流编辑器
    awk - 文本分析工具
    clang - 编译工具
    netstat - 查看网络状态，包含路由表信息
    xcodebuild - 编译xcode工程
    traceroute - 查看路由路线
    lsof - 查看端口使用情况    
    rvictl - iOS设备网卡映射
    afplay - 播放本地音乐
    mplayer - 播放音视频
    system_profiler - 查看系统软硬件配置
    stunclient - 获取公网IP
    screencapture - 屏幕抓拍
    open - 打开文件
    lipo - 查看可执行文件CPU架构
    pbcopy - 复制到剪贴板
    pbpaste - 黏贴
    say - 播放文字
    diskutil - 磁盘工具
    wc - 统计行数
    nslookup - 根据域名查看ip
    osascript - 执行AppleScript、JavaScript
    设置音量
    osascript -e 'set volume 1'
    networksetup - 网络设置
    连接WIFI
    networksetup -setairportnetwork en0 WIFI_SSID WIFI_PASSWORD
    //Enable
    dsenableroot
    //Disable
    dsenableroot -d


参考链接:
1. [linux 创建连接命令 ln -s 软链接](http://www.cnblogs.com/kex1n/p/5193826.html)
1. [centos 7 查看内网ip和外网ip](http://blog.csdn.net/qq_31382921/article/details/53836523)
2. [Mac转发80端口](https://www.jianshu.com/p/58ec8f1e480d)
3. [linux 创建连接命令 ln -s 软链接](http://www.cnblogs.com/kex1n/p/5193826.html)
4. [Mac下配置Redis服务器（自启动、后台运行）](http://blog.csdn.net/langzi7758521/article/details/51684413)
5. [iTerm2使用rz、sz远程上传或下载文件](https://www.noonme.com/post/2016/02/mac-iterm2-rz-sz/)
6. [Example syntax for Secure Copy (scp)](http://www.hypexr.org/linux_scp_help.php)
7. [linux下查看文件编码及修改编码](https://blog.csdn.net/jnbbwyth/article/details/6991425)
8. [Linux通过PID查看进程完整信息](https://blog.csdn.net/great_smile/article/details/50114133)
9. [工具参考篇](https://linuxtools-rst.readthedocs.io/zh_CN/latest/tool/crontab.html)
