.PHONY: all

build:
	docker build -t ${tag} .

push:
	docker tag ${tag} knabben/aws-tools:${tag}
	docker push knabben/aws-tools:${tag}
