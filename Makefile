IMAGE="your image name"
TAG="dev"

BRANCH=$(shell git symbolic-ref --short -q HEAD)
COMMIT=$(shell git rev-parse --short HEAD)

release: image image.push

.PHONY: build
build:
	swag init && go build -o ./tmp/main .

.PHONY: image
image:
	make build
	docker build --build-arg branch=${BRANCH} --build-arg commit=${COMMIT} --tag ${IMAGE}:${TAG} .

.PHONY: image.push
image.push:
	docker push ${IMAGE}:${TAG}
