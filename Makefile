GO_BUILD_OS ?= linux
GO_BUILD_OUTPUT ?= main
GO_SERVICE_NAME ?= go-whatsapp-rest
GO_EXPOSE_PORT ?= 3000
DOCKER_IMAGE_NAME ?= go-whatsapp-rest
DOCKER_IMAGE_VERSION ?= latest

git-push:
	make go-dep-init
	make go-dep-clean
	make go-clean
	git add .
	git commit -am "$(COMMIT_MSG)"
	git push origin master

git-pull:
	git pull origin master

go-dep-init:
	make go-dep-clean
	rm -f Gopkg.toml Gopkg.lock
	dep init -v

go-dep-ensure:
	make go-dep-clean
	dep ensure -v

go-dep-clean:
	rm -rf ./vendor

go-build:
	make go-clean
	make go-dep-ensure
	CGO_ENABLED=0 GOOS=$(GO_BUILD_OS) go build -a -o ./build/$(GO_BUILD_OUTPUT) *.go

go-run:
	CONFIG_ENV="DEV" CONFIG_FILE_PATH="./build/configs" CONFIG_LOG_LEVEL="DEBUG" CONFIG_LOG_SERVICE="$(GO_SERVICE_NAME)" go run *.go

go-clean:
	rm -f ./build/$(GO_BUILD_OUTPUT)

docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) --build-arg SERVICE_NAME=$(GO_SERVICE_NAME) .

docker-run:
	docker run -d -p $(GO_EXPOSE_PORT):$(GO_EXPOSE_PORT) --name $(GO_SERVICE_NAME) --rm $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION)
	make docker-logs

docker-exec:
	docker exec -it $(GO_SERVICE_NAME) bash

docker-stop:
	docker stop $(GO_SERVICE_NAME)

docker-logs:
	docker logs $(GO_SERVICE_NAME)

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION)

docker-clean:
	docker rmi -f $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION)

clean:
	make go-clean
	make go-dep-clean
	make docker-clean
