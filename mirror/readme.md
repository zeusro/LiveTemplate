收集各种镜像资源。

[阿里巴巴开源镜像](https://opsx.alibaba.com/mirror?lang=zh-CN)


## alpine修改镜像源


使用阿里镜像 https://mirrors.aliyun.com


  sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

使用科大镜像 http://mirrors.ustc.edu.cn

  sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

[alpine修改镜像源](https://www.jianshu.com/p/791c91b7c2a4)
