
# CGO_ENABLED=0 because of https://stackoverflow.com/a/55106860/1716310

BINARY_NAME=pi_monitor_api
VERSION=$(shell cat ./version.txt)
BINARY_PORT?=8080

IMAGE_NAME?=pi_monitor_api
IMAGE_VERSION=$(shell cat ./version.txt)

CONTAINER_PORT?=8080
API_CONTAINER_ALIAS?=pi_monitor

# Enable experimental features for Docker to build for other archs
export DOCKER_CLI_EXPERIMENTAL=enabled

build:
	cd ./src; \
	go mod tidy; \
	CGO_ENABLED=0 go build -o ${BINARY_NAME}_${VERSION} cmd/main.go

build-linux-amd64: delete-builds
	cd ./src; \
	go mod tidy; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}_${VERSION}_amd64 cmd/main.go

build-linux-armv7: delete-builds
	cd ./src; \
	go mod tidy; \
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ${BINARY_NAME}_${VERSION}_armv7 cmd/main.go

build-linux-arm64: delete-builds
	cd ./src; \
	go mod tidy; \
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}_${VERSION}_arm64 cmd/main.go

delete-builds:
	-cd ./src; \
	rm -f ${BINARY_NAME}_*

create-builder:
	if ! docker buildx ls | grep -q 'multiplatform'; then \
		docker buildx create --name multiplatform --use; \
	fi

run:
	mkdir -p ./logs; \
	./src/${BINARY_NAME}_${VERSION} > ./logs/${BINARY_NAME}_${VERSION}.log 2>&1 &

test:
	cd ./src; \
	go mod tidy; \
	go test ./... -cover -count=1

local-run-debug:
	cd ./src; \
	go mod tidy; \
	go run cmd/main.go

logs:
	tail -f ./logs/${BINARY_NAME}_${VERSION}.log

build-docker-image: build
	docker build --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:${VERSION} .

build-docker-image-amd64: build-linux-amd64
	docker buildx build --load --platform linux/amd64 --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:${IMAGE_VERSION} .

build-docker-image-armv7: build-linux-armv7
	docker buildx build --load --platform linux/arm/v7 --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:${IMAGE_VERSION} .

build-docker-image-arm64: build-linux-arm64
	docker buildx build --load --platform linux/arm64 --build-arg VERSION=${VERSION} -t ${IMAGE_NAME}:${IMAGE_VERSION} .

run-docker-image:
	docker run \
	--name ${API_CONTAINER_ALIAS} \
	-d \
	-p ${CONTAINER_PORT}:${BINARY_PORT} \
	--security-opt systempaths=unconfined \
	${IMAGE_NAME}:${IMAGE_VERSION}

stop:
	@echo "Stopping Docker container..."
	-docker stop $(CONTAINER_NAME)

container-rm:
	@echo "Deleting Docker container..."
	-docker rm -f ${API_CONTAINER_ALIAS}

container-debug:
	docker exec -it ${API_CONTAINER_ALIAS} /bin/sh

container-logs:
	docker logs -f `docker ps -qa --filter name=${API_CONTAINER_ALIAS}`

container-stats:
	docker stats `docker ps -q --filter name=${API_CONTAINER_ALIAS}`

container-redeploy: stop container-rm delete-builds build build-docker-image run-docker-image

gittag:
	git tag -a ${IMAGE_VERSION} -m "Release ${IMAGE_VERSION}"
	git push origin ${IMAGE_VERSION}
	@echo "Updating version..."
	@echo "$(shell cat version.txt)-SNAPSHOT" > version.txt
	@echo "Version updated to $(shell cat version.txt)"

gittags:
	git tag --sort version:refname | tail
