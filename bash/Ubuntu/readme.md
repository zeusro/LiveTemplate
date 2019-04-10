
1. [安装Java8](#安装Java8)
1. [安装 Java10](#安装Java10)
1. [获取](#获取)
1. [颜色标记](#颜色标记)
1. [18安装docker](#18安装docker)
1. [18安装docker-compose](#18安装docker-compose)
1. []()

## 18安装docker-compose
```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## 18安装docker

```bash
sudo apt-get remove -y docker docker-engine docker.io
sudo apt-get -y update
sudo apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get -y update
sudo apt-get install -y docker-ce
```


## 安装 Java8

```bash
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get update
sudo apt-get install oracle-java8-installer
sudo apt install oracle-java8-set-default
```

## Ubuntu16 安装 Java8

```bash
sudo apt-get update
sudo apt-get install -y default-jre
sudo apt-get install -y default-jdk

```

## 安装 Java10

```bash
install_java10(){
  # 参考https://stackoverflow.com/questions/49507160/how-to-install-jdk-10-under-ubuntu
  wget https://download.java.net/openjdk/jdk10/ri/jdk-10_linux-x64_bin_ri.tar.gz  
  sudo tar -zxvf jdk-10_linux-x64_bin_ri.tar.gz -C /usr/lib
  sudo update-alternatives --install /usr/bin/java java /usr/lib/jdk-10/bin/java 1
  sudo update-alternatives --install /usr/bin/javac javac /usr/lib/jdk-10/bin/javac 1  
  sudo update-alternatives --config java
  # sudo apt-get install oracle-java10-installer
}
```

## 获取

```bash
sysArch(){
    ARCH=$(uname -m)
    if [[ "$ARCH" == "i686" ]] || [[ "$ARCH" == "i386" ]]; then
        VDIS="32"
    elif [[ "$ARCH" == *"armv7"* ]] || [[ "$ARCH" == "armv6l" ]]; then
        VDIS="arm"
    elif [[ "$ARCH" == *"armv8"* ]] || [[ "$ARCH" == "aarch64" ]]; then
        VDIS="arm64"
    elif [[ "$ARCH" == *"mips64le"* ]]; then
        VDIS="mips64le"
    elif [[ "$ARCH" == *"mips64"* ]]; then
        VDIS="mips64"
    elif [[ "$ARCH" == *"mipsle"* ]]; then
        VDIS="mipsle"
    elif [[ "$ARCH" == *"mips"* ]]; then
        VDIS="mips"
    elif [[ "$ARCH" == *"s390x"* ]]; then
        VDIS="s390x"
    fi
    return 0
}
```

## 颜色标记

```bash

colorEcho(){
    COLOR=$1
    echo -e "\033[${COLOR}${@:2}\033[0m"
}
```

## 安装redis

```bash
sudo apt-get install -y build-essential tcl
cd /tmp
curl -O http://download.redis.io/redis-stable.tar.gz
tar xzvf redis-stable.tar.gz
cd redis-stable
make
make test
sudo make install
```

配置
```bash
sudo mkdir /etc/redis
sudo cp /tmp/redis-stable/redis.conf /etc/redis
sudo nano /etc/redis/redis.conf
```

其余具体见参考链接...

需要注意的是
`sudo systemctl start redis`改为`sudo systemctl unmask  redis-server.service`

sudo systemctl unmask status redis

## 参考链接:

1. [Ubuntu安装Java8和Java9](https://www.cnblogs.com/woshimrf/p/ubuntu-install-java.html)
1. [BASH的基本语法](https://www.cnblogs.com/lonelywolfmoutain/p/5950439.html)
1. [How To Install and Configure Redis on Ubuntu 16.04](https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-redis-on-ubuntu-16-04)
1. [Failed to start redis.service: Unit redis-server.service is masked
Ask](https://stackoverflow.com/questions/40317106/failed-to-start-redis-service-unit-redis-server-service-is-masked)
1. []()
1. []()
1. []()