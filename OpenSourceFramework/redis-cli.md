
```
172.18.221.62:31303

redis-cli -h host -p port -a password

redis-cli -h 172.18.221.62 -p 31303 -a password --raw


redis-cli -h 172.31.7.44 -p 6379  --raw

# redis 5.0
redis-cli -h 172.18.221.62 -p 6379  --raw

redis-benchmark  -h 172.18.221.62 -p 6379  -q -n 100000

```

kex $po sh
redis-cli -h codis-server-0 -p 6379  --raw

## bigKeys

redis-cli -h $host -p $p --bigkeys

    redis-cli -h $host  --bigkeys


# 客户端列表
client list
# 获取10个慢日志
SLOWLOG GET 10
# 遍历key
keys market-bq:*


redis-cli -h codis-server-0.codis-server.17zwd.svc.cluster.local -p 6379  --raw

redis-cli -h codis-server-12.codis-server.17zwd.svc.cluster.local -p 6379  --raw

redis-cli -h 172.31.1.144 -p 6379  --raw



# 读写


```

EXPIRE foo 1
EXPIRE foo <秒>
TTL
```





## 终极命令

FLUSHALL

## 参考

[redis cli命令](https://www.cnblogs.com/kongzhongqijing/p/6867960.html)


https://redis.io/topics/rediscli

https://blog.csdn.net/yangcs2009/article/details/50781530

https://www.runoob.com/redis/redis-benchmarks.html

https://my.oschina.net/davehe/blog/174662
