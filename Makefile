VERSION:=$(shell git rev-parse HEAD)

start:
	go run ./server/main.go

exp-docker-cloud:
	docker build -t us.gcr.io/constant-bolt-126016/stl-com:latest .
	docker push us.gcr.io/constant-bolt-126016/stl-com:latest

.PHONY: start exp-docker-cloud