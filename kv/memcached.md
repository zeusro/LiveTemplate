# memcached

## 监听的端口
```
-p 
```
## 连接的IP地址, 默认是本机
```
-l 
```

## 启动memcached服务
```
-d start
```

## 重起memcached服务
```
-d restart 
```


## 关闭正在运行的memcached服务
```
-d stop|shutdown 
```


## 安装memcached服务
```
-d install
```

## 卸载memcached服务
```
-d uninstall 
```


## 以管理员的身份运行 (仅在以root运行的时候有效)
```
-u 
```


## 最大内存使用，单位MB。默认64MB
```
-m 
```

## 内存耗尽时返回错误，而不是删除项
```
-M 
```

## 最大同时连接数，默认是1024
```
-c 
```

## 块大小增长因子，默认是1.25
```
-f 
```

## 最小分配空间，key+value+flags默认是48
```
-n 
```

## 显示帮助
```
-h 
```