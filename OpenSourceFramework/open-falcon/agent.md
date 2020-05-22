

## 一键安装

```
wget --no-check-certificate -O open-falcon-agent.sh https://gist.githubusercontent.com/zeusro/f42695c2e9044cb6a009f9cff50c996b/raw/6c8d431304cbe3a53be6b8a91715fe0dc4c43131/open-falcon-agent.sh
bash open-falcon-agent.sh
# bash /open-falcon-agent.sh > agent.log

```

## 初始配置
```
sed -i 's/<old>/<new>/g'  $workPath/$appName/agent/config/cfg.json

export workPath="/zeusro"
export appName="falcon-plus"
sed -i 's/0\.0\.0\.0\:6030/120\.78\.148\.177\:6030/g'  $workPath/$appName/agent/config/cfg.json
sed -i 's/0\.0\.0\.0\:8433/120\.78\.148.177\:8433/g'  $workPath/$appName/agent/config/cfg.json
sed -i 's/debug\"\:[ ]*true/debug\"\: false/g'  $workPath/$appName/agent/config/cfg.json

vi $workPath/$appName/agent/config/cfg.json
cd $workPath/$appName
./open-falcon start agent
./open-falcon restart agent
./open-falcon check agent


```

## 检查

```

./open-falcon check
./open-falcon restart alarm

vi $workPath/$appName/alarm/config/cfg.json
cat $workPath/$appName/alarm/log/cfg.json

cat /zeusro/falcon-plus/alarm/logs/alarm.log

 sed -i 's/http:\/\/127.0.0.1:10086\/sms/newsms/g' $workPath/$appName/alarm/config/cfg.json
 
 sed -i 's/newsms/http:\/\/vlog\.17zwd\.com\/notify\/sms?token=token/g' $workPath/$appName/alarm/config/cfg.json
 
 cat $workPath/$appName/alarm/config/cfg.json
```

## 检查日志

```
cat $workPath/$appName/agent/logs/agent.log
cat $workPath/$appName/alarm/logs/alarm.log 
```

```

root@iZwz9b1zq3z6108bj4u84tZ:~# cat

cat /zeusro/falcon-plus/alarm/logs/alarm.log 
/zeusro/falcon-plus/alarm/logs/alarm.log | grep vlog
time="2018-03-21T11:41:19+08:00" level=debug msg="send sms:<Tos:18011701109, Content:[P0][PROBLEM][120.79.9.50-web-pdd][][ all(#3) proc.num name=php-fpm 0==0][O3 2018-03-21 11:41:00]>, resp:[{\"code\":48,\"msg\":\"模板参数总长度超过限制\",\"detail\":\"模板中使用的参数总长度为108，请勿超过10\",\"data\":null}], url:http://vlog.17zwd.com/notify/sms?token=token"   
```

## 重启

```
cd $workPath/$appName
rm -rf $workPath/$appName/agent/logs/agent.log
./open-falcon restart agent 
```

```
```








