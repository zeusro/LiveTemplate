

```
cd "C:\Program Files\Redis"
redis-server.exe  redis.windows.conf
.\redis-server.exe .\redis.windows.conf
.\redis-server --service-start .\redis.windows.conf
.\redis-server --service-install redis.windows.conf --loglevel verbose
.\redis-server --service-install --service-name redisService1 --port 10001
redis-server --service-start --service-name redisService1
redis-server --service-uninstall --service-name redisService1


```

[Running Redis as a Service](https://github.com/MicrosoftArchive/redis/blob/3.0/Windows%20Service%20Documentation.md)