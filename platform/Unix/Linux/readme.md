
## 文件类

* du -h --max-depth=1

查看目录占用了多少空间

* ll
 
ls -a 的别称

* 找完删除

find . -name "*.log"  | xargs rm -rf

```
# 统计文件个数
ls -l |grep "^-"|wc -l

# 统计文件夹个数
ls -l |grep "^ｄ"|wc -l
```

## 其他

* 查看系统版本

```bash
cat /proc/version
cat /etc/redhat-release
cat /etc/issue
```

* 将命令的结果定义为变量
```bash
var=$(command)
```
* 一键脚本
```bash
wget --no-check-certificate -O shadowsocks-all.sh https://raw.githubusercontent.com/teddysun/shadowsocks_install/master/shadowsocks-all.sh
chmod +x shadowsocks-all.sh
bash ./shadowsocks-all.sh 2>&1 | tee shadowsocks-all.log
```


* cp
```
cp -r /public /mnt/public
```

```
cp [options] <source file or directory> <target file or directory>

或

cp [options] source1 source2 source3 …. directory

上面第一条命令为单个文件或目录拷贝，下一个为多个文件拷贝到最后的目录。

options选项包括：

- a 保留链接和文件属性，递归拷贝目录，相当于下面的d、p、r三个选项组合。
- d 拷贝时保留链接。
- f 删除已经存在目标文件而不提示。
- i 覆盖目标文件前将给出确认提示，属交互式拷贝。
- p 复制源文件内容后，还将把其修改时间和访问权限也复制到新文件中。
- r 若源文件是一目录文件，此时cp将递归复制该目录下所有的子目录和文件。当然，目标文件必须为一个目录名。
- l 不作拷贝，只是链接文件。
-s 复制成符号连结文件 (symbolic link)，亦即『快捷方式』档案；
-u 若 destination 比 source 旧才更新 destination。
```

* nohup

    nohup 命令 >/dev/null 2>&1 &

参考链接:
1. [find](https://blog.csdn.net/ydfok/article/details/1486451)
2. [cp](http://www.metsky.com/archives/542.html)
3. 

