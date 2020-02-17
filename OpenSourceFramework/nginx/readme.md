

## 安装

* Ubuntu
```
apt-get -y update
apt-get -y install nginx
nginx -v
nginx -V
```

* Centos
```

```

## 配置位置


* Ubuntu

```
vi /etc/nginx/nginx.conf
```

* mac

```
which nginx
cd /usr/local/Cellar/nginx/1.13.9

sudo vi /usr/local/etc/nginx/nginx.conf
```
## 静态资源位置

* mac-brew

```
/usr/local/var/www/
```

sudo chmod -R 777  /usr/local/var/www/pdd/resources

## 重启

* Ubuntu

```
nginx -s reload;
```

* mac

```
cd /usr/local/nginx/sbin/;
./nginx -s reload;
```


## 常用配置


* 设置首页
```
    server {
        listen 80;
        server_name www.angelina.ink;
        root  /angelina/birth;
        location / {
        index index.html;
                }
        } 
```

* 反向代理
```
    server {
        listen       80;
        server_name  note.zeusro.tech;

        location / {
            proxy_pass http://localhost:9000;
          } 
          }
```


## 配置相关

### ssl

```
server {
        listen 443 ssl http2;
        server_name example.com www.example.com;

        # SSL/TLS configs
        ssl off;
        ssl_certificate /etc/ssl/certs/example_com_cert_chain.crt;
        ssl_certificate_key /etc/ssl/private/example_com.key;


}
```

模式|含义
----|----
location ^~|uri以某个常规字符串开头，正则之前
location ~|区分大小写的正则匹配;
location  ~*|不区分大小写的正则匹配
location =|精确匹配,只有完全匹配才生效
location /uri|同^~ ,正则之后
location /|通用匹配


### 配置校验


    nginx -t -c /usr/local/nginx/conf/nginx.conf
    nginx -t -c /etc/nginx/nginx.conf

### 重定向 

```
server {
  server_name .mydomain.com;
  return 301 http://www.adifferentdomain.com$request_uri;
}
```

- http->https
```
server{
        listen 80;
        server_name example.com www.example.com;
        return 301 https://www.example.com$request_uri;
}
```

    
## 参考链接
1. [多条件if](https://gist.github.com/jrom/1760790)
1. https://www.nginx.com/resources/wiki/start/topics/examples/full/
1. https://nginx.org/en/docs/
1. https://docs.nginx.com/nginx/admin-guide/
2. [nginx 基础配置：多个location转发任意请求或访问静态资源文件](https://blog.csdn.net/tutian2000/article/details/81531513)
3. [环境变量](http://nginx.org/en/docs/http/ngx_http_core_module.html)
1. [nginx 启动错误"nginx: [emerg] host not found in upstream "解决方案](https://blog.csdn.net/WangXiaoMing099/article/details/23443623)
1. [在容器中利用Nginx-proxy实现多域名的自动反向代理、免费SSL证书](https://www.jianshu.com/p/2c6fd9e43aa7)
2. [doc](https://docs.nginx.com/)
