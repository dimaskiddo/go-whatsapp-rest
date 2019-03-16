BUILD_OS           := linux
BUILD_OUTPUT       := main
SERVICE_NAME       := go-whatsapp-rest
SERVICE_PORT       := 3000
IMAGE_NAME         := go-whatsapp-rest
IMAGE_TAG          := latest
REBASE_URL         := "github.com/dimaskiddo/go-whatsapp-rest"
COMMIT_MSG         := "update improvement"

.PHONY:

.SILENT:

init:
	make clean
	rm -f Gopkg.toml Gopkg.lock
	dep init -v

ensure:
	make clean
	dep ensure -v

build:
	make clean
	make ensure
	CGO_ENABLED=0 GOOS=$(BUILD_OS) go build -a -o ./build/$(BUILD_OUTPUT) *.go
	echo "Build complete please check build directory."

run:
	CONFIG_ENV="DEV" CONFIG_FILE_PATH="./build/configs" CONFIG_LOG_LEVEL="DEBUG" CONFIG_LOG_SERVICE="$(SERVICE_NAME)" go run *.go

clean:
	rm -f ./build/$(BUILD_OUTPUT)
	rm -rf ./vendor

commit:
	make init
	make clean
	git add .
	git commit -am "$(COMMIT_MSG)"

rebase:
	rm -rf .git
	find . -type f -iname "*.go*" -exec sed -i '' -e "s%github.com/dimaskiddo/go-whatsapp-rest%$(REBASE_URL)%g" {} \;	
	git init
	git remote add origin https://$(REBASE_URL)

push:
	git push origin master

pull:
	git pull origin master

c-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) --build-arg SERVICE_NAME=$(SERVICE_NAME) .

c-run:
	docker run -d -p $(SERVICE_PORT):$(SERVICE_PORT) --name $(SERVICE_NAME) --rm $(IMAGE_NAME):$(IMAGE_TAG)
	make docker-logs

c-shell:
	docker exec -it $(SERVICE_NAME) bash

c-stop:
	docker stop $(SERVICE_NAME)

c-logs:
	docker logs $(SERVICE_NAME)

c-push:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

c-clean:
	docker rmi -f $(IMAGE_NAME):$(IMAGE_TAG)
