## docker 配置位置

     ~/.docker/daemon.json
    
## bin

    /usr/local/bin

## maven 地址

```bash
/usr/local/Cellar/maven/3.5.2

settings.xml
/usr/local/Cellar/maven/3.5.2/libexec/conf/settings.xml
$usr/.m2/repository
/Volumes/D/hy
```

* localrepository

```
/Users/zeusro/.m2/repository
```

## JAVA_HOME

```bash
/Library/Java/JavaVirtualMachines/jdk1.8.0_151.jdk/Contents/Home
/Library/Java/JavaVirtualMachines/jdk1.8.0_151.jdk/Contents/Home/jre
export JAVA_HOME=`/usr/libexec/java_home -v 1.8` 
echo  $JAVA_HOME
```

## Nginx

/usr/local/etc/nginx/nginx.conf

## Go

* goPath

go env
/Users/zeusro/go

* goRoot

 /usr/local/Cellar/go

## host

/etc/hosts

## path

1. /etc/profile   （建议不修改这个文件 ）

 全局（公有）配置，不管是哪个用户，登录时都会读取该文件。

2. /etc/bashrc    （一般在这个文件中添加系统级环境变量）

 全局（公有）配置，bash shell执行时，不管是何种方式，都会读取此文件。

3. ~/.bash_profile  （一般改这个）

/usr/local/bin 放在 /usr/bin 前面

 每个用户都可使用该文件输入专用于自己使用的shell信息,当用户登录时,该文件仅仅执行一次!

4. /private/etc/paths

## zsh

~/.zshrc

##  docker 位置

~/.docker/daemon.json
