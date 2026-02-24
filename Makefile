GO = go
PROJECT = cine_stream
VERSION = $(shell date +%m%d%H%M)
GITLAB = registry.gitlab.com/cinemae
URL = $(GITLAB)/$(PROJECT)
REGISTRY_RELEASE = $(URL):$(VERSION)
REGISTRY_LATEST = $(URL):latest

$(PROJECT):
	$(GO) build -o $@ .

net:
	@docker network inspect shared-network >/dev/null 2>&1 || docker network create shared-network

test:
	$(GO) test $(M)  -v -gcflags=all=-l -coverpkg=./... -coverprofile=test.out ./...

clean:
	rm -f $(PROJECT)

linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(PROJECT)

run:
	go run main.go --conf=conf/dev/app.yaml

up:
	docker compose up -d gate stream --scale mysql=0 --scale redis=0

up-all:
	docker compose up -d

down:
	docker compose down

migrate:
	go run main.go --conf=conf/dev/app.yaml --migrate

build:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(PROJECT)
	sh ./build.sh
	docker build -t $(PROJECT):latest .
	rm -fr ./release

push:
	docker tag $(PROJECT) $(REGISTRY_RELEASE)
	docker tag $(PROJECT) $(REGISTRY_LATEST)
	docker push $(REGISTRY_RELEASE)
	docker push $(REGISTRY_LATEST)
	docker rmi $(REGISTRY_RELEASE)
	docker rmi $(REGISTRY_LATEST)