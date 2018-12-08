


- 相关文档

https://www.consul.io/api/kv.html


- 获取配置

/v1/kv/goms/config/17.java-spring-boot-template


- 新增/修改配置

```
put
/v1/kv/goms/config/17.java-spring-boot-template2
```

```
contents=`cat /Volumes/D/a.conf`
curl \
    --request PUT \
    --data-binary $contents \
    /v1/kv/goms/config/17.java-spring-boot-template2
```
