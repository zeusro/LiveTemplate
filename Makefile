# 时区偏移模式
now := $(shell date +"%Y-%m-%d %H:%M:%S %z" | sed 's/\([+-][0-9][0-9]\)\([0-9][0-9]\)/\1:\2/')

auto_commit:
	git add .
	git commit -am "Update: $(now)"
	git pull
	git push