.PHONY: all

compile:
	go build main.go

build:
	docker build -t aws-tools:${tag} .

push:
	docker tag ${tag} knabben/aws-tools:${tag}
	docker push knabben/aws-tools:${tag}
