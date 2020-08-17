
## 下载最新版

```bash
curl -s $releases_page | grep "browser_download_url" | cut -d '"' -f 4


releases_page='https://api.github.com/repos/alibaba/kt-connect/releases/latest'

curl -s $releases_page \
	| grep "browser_download_url" \
	| cut -d '"' -f 4 \
	| fzf \
	| curl -O


curl -s $releases_page | grep "browser_download_url" | cut -d '"' -f 4
```
