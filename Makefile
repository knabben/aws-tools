.PHONY: all

build:
	docker build -t ${tag} .
