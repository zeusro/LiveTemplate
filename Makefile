now := $(shell date +"%Y-%m-%d %H:%M:%S %Z")

auto_commit:
	git add .
	git commit -am "Update: $(now)"
	git pull
	git push