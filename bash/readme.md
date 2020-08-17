
## 获取IP

    curl https://www.ipip.net/ip.html | grep name\=\"ip\"

访问多了需要这样:

    curl https://www.ipip.net/ip.html | grep  iframe

## 网络诊断

    iftop -i eth0 -nNB -m 10M