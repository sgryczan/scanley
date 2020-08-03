VERSION=0.0.1
IMAGE_NAME=sgryczan/scanley
IMAGE_TAG=${VERSION}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the image
	docker build --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:${IMAGE_TAG} .

.PHONY: push
push: ## Push the image
	docker push ${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: run
run: ## Run the image
	docker run -p 8080:8080 ${IMAGE_NAME}:${IMAGE_TAG}


.PHONY: test
test: ## Run tests
	go test -v ./... -cover
