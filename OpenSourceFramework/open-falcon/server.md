
```bash

sudo apt-get install -y redis-server  mysql-server git

getconf LONG_BIT

export workPath="/zeusro"
export appName="falcon-plus"

echo $appName
```

# server


```bash
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/falcon-plus.git

cd $GOPATH/src/github.com/open-falcon/falcon-plus/scripts/mysql/db_schema/
mysql -h 127.0.0.1 -u root -pROOT < 1_uic-db-schema.sql
mysql -h 127.0.0.1 -u root -pROOT < 2_portal-db-schema.sql
mysql -h 127.0.0.1 -u root -pROOT < 3_dashboard-db-schema.sql
mysql -h 127.0.0.1 -u root -pROOT < 4_graph-db-schema.sql
mysql -h 127.0.0.1 -u root -pROOT < 5_alarms-db-schema.sql

cd $GOPATH/src/github.com/open-falcon/falcon-plus/
make all
#make 失败的话运行以下命令
sudo apt-get install --reinstall make

make pack

mkdir -p $workPath/$appName
tar -xzvf open-falcon-v0.2.1.tar.gz -C $workPath/$appName


cd $workPath/$appName
cat $workPath/$appName/api/config/cfg.json
cat $workPath/$appName/aggregator/config/cfg.json
cat $workPath/$appName/alarm/config/cfg.json
cat $workPath/$appName/graph/config/cfg.json
cat $workPath/$appName/nodata/config/cfg.json
cat $workPath/$appName/hbs/config/cfg.json

grep -Ilr 3306  ./ | xargs -n1 -- sed -i 's/\"root:/\"root:SwitchBest3!$/g'


grep -Ilr 3306  ./ | xargs -n1 -- sed -i 's/SwitchBest3\!\$/20180315SwitchBest/g'

grep -Ilr 3306  ./ | xargs -n1 -- sed -i 's/\"root:/\"root:<密码>/g'


cat $workPath/$appName/api/logs/api.log

nohup /jp/open-falcon start > /dev/null 2>&1 &

./open-falcon start

./open-falcon start api
```

* 配置

```bash
vi $workPath/$appName/api/config/cfg.json
vi $workPath/$appName/aggregator/config/cfg.json
vi $workPath/$appName/alarm/config/cfg.json
vi $workPath/$appName/graph/config/cfg.json
vi $workPath/$appName/nodata/config/cfg.json
vi $workPath/$appName/hbs/config/cfg.json
```

# 安装仪表盘

* docker 模式

```bash
sudo apt-get install -y python-virtualenv slapd ldap-utils libmysqld-dev build-essential python-dev libldap2-dev docker.io


git clone https://github.com/open-falcon/dashboard.git
cd $workPath/dashboard
# 记得先用阿里云镜像给docker 加速
docker build -t falcon-dashboard:v1.0 .

docker run -itd --name aaa --net host \
	-e API_ADDR=http://127.0.0.1:8080/api/v1 \
	-e PORTAL_DB_HOST=127.0.0.1 \
	-e PORTAL_DB_PORT=3306 \
	-e PORTAL_DB_USER=root \
	-e PORTAL_DB_PASS=123456 \
	-e PORTAL_DB_NAME=falcon_portal \
	-e ALARM_DB_PASS=123456 \
	-e ALARM_DB_HOST=127.0.0.1 \
	-e ALARM_DB_PORT=3306 \
	-e ALARM_DB_USER=root \
	-e ALARM_DB_PASS=123456 \
	-e ALARM_DB_NAME=alarms \
	falcon-dashboard:v1.0
```

```
docker run -itd --name aaa --net host \
	-e API_ADDR=http://127.0.0.1:8080/api/v1 \
	-e PORTAL_DB_HOST=127.0.0.1 \
	-e PORTAL_DB_PORT=3306 \
	-e PORTAL_DB_USER=root \
	-e PORTAL_DB_PASS=20180315SwitchBest \
	-e PORTAL_DB_NAME=falcon_portal \
	-e ALARM_DB_PASS=20180315SwitchBest \
	-e ALARM_DB_HOST=127.0.0.1 \
	-e ALARM_DB_PORT=3306 \
	-e ALARM_DB_USER=root \
	-e ALARM_DB_PASS=20180315SwitchBest \
	-e ALARM_DB_NAME=alarms \
	falcon-dashboard:v1.0
```

## 第二次启动

```bash
docker ps -a 
docker restart e4ec31872dda --net host \
	-e API_ADDR=http://127.0.0.1:8080/api/v1 \
	-e PORTAL_DB_HOST=127.0.0.1 \
	-e PORTAL_DB_PORT=3306 \
	-e PORTAL_DB_USER=root \
	-e PORTAL_DB_PASS=20180315SwitchBest \
	-e PORTAL_DB_NAME=falcon_portal \
	-e ALARM_DB_PASS=20180315SwitchBest \
	-e ALARM_DB_HOST=127.0.0.1 \
	-e ALARM_DB_PORT=3306 \
	-e ALARM_DB_USER=root \
	-e ALARM_DB_PASS=20180315SwitchBest \
	-e ALARM_DB_NAME=alarms \
	falcon-dashboard:v1.0
```

telnet 120.78.146.177 8081

http://127.0.0.1:8081/

* 普通模式

```bash
sudo apt-get install -y python-virtualenv slapd ldap-utils libmysqld-dev build-essential python-dev libldap2-dev

sudo apt-get install -y python-setuptools  
easy_install virtualenv 


cd $workPath/dashboard

mkdir -p $workPath/$appName

pip install -U setuptools
./env/bin/pip install -r pip_requirements.txt -i https://pypi.douban.com/simple

rm -rf /zeusro/dashboard
sudo ./env/bin/python wsgi.py

index-url = https://pypi.mirrors.ustc.edu.cn/simple/ 

trusted-host=pypi.mirrors.ustc.edu.cn
./env/bin/pip install -r pip_requirements.txt -i https://pypi.mirrors.ustc.edu.cn/simple/  

cd /zeusro/dashboard
sudo bash control stop
vi /zeusro/dashboard/rrd/config.py
cat /zeusro/dashboard/rrd/config.py
sudo vi /zeusro/dashboard/rrd/config.py

sudo ufw allow proto tcp to 0.0.0.0/0 port 8081


/usr/local/go/src/github.com/open-falcon


grep -Ilr 3306  ./ | sudo xargs -n1 -- sed -i 's/root:password/root:ROOT/g'

```