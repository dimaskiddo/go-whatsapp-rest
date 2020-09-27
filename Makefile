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
	GO111MODULE=on go mod init

init-dist:
	mkdir -p dist
	touch dist/.gitkeep

vendor:
	make clean
	GO111MODULE=on go mod vendor

release:
	make vendor
	goreleaser --snapshot --skip-publish --rm-dist
	make init-dist
	echo "Release complete please check dist directory."

publish:
	GITHUB_TOKEN=$(GITHUB_TOKEN) goreleaser --rm-dist
	make clean-dist
	echo "Publish complete please check your repository releases."

run:
	go run *.go

clean-dist:
	rm -rf ./dist/*
	make init-dist

clean:
	make clean-dist
	rm -rf ./vendor

commit:
	make vendor
	make clean
	git add .
	git commit -am "$(COMMIT_MSG)"

rebase:
	rm -rf .git
	find . -type f -iname "*.go*" -exec sed -i '' -e "s%github.com/dimaskiddo/go-whatsapp-rest%$(REBASE_URL)%g" {} \;
	git init
	git remote add origin https://$(REBASE_URL).git

push:
	git push origin master

pull:
	git pull origin master

c-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) --build-arg SERVICE_NAME=$(SERVICE_NAME) .

c-run:
	docker run -d -p $(SERVICE_PORT):$(SERVICE_PORT) --name $(SERVICE_NAME) --rm $(IMAGE_NAME):$(IMAGE_TAG)
	make c-logs

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
