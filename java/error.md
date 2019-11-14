
## mysql ssl 警告

```
Establishing SSL connection without server's identity verification is not recommended. According to MySQL 5.5.45+, 5.6.26+ and 5.7.6+ requirements SSL connection must be established by default if explicit option isn't set. For compliance with existing applications not using SSL the verifyServerCertificate property is set to 'false'. You need either to explicitly disable SSL by setting useSSL=false, or set useSSL=true and provide truststore for server certificate verification.
```

连接字符串加上 useSSL=false/true 即可 


## mvn 3.6.2

不要用,构建会有问题

## Gradle 6.0 多模块构建失败

出现这个问题我也很懵逼A->B->C

A用了C的方法,但是提示报错,最后我只能在A的build.gradle里面显示依赖C,这个bug很恶心.
