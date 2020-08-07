now := $(shell date)

auto_commit: pull
	git add .
	git commit -am "$(now)"	
	git push

gitbook:
	pip install sphinx recommonmark
	sphinx-quickstart
	pip freeze > requirements.txt
