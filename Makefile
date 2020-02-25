BINARY_NAME=socks5
DOCKER_IMAGE=gosocks5
DOCKER_CONTAINER=socks5
VERSION=0.1

all: clean deps test build
build:
	go build -o ${BINARY_NAME} cmd/main.go
test:
	go test -v -race
clean:
	go clean
	rm -f ${BINARY_NAME}
deps:
	go mod download
run: build
	./socks5

docker-build:
	docker build -t ${DOCKER_IMAGE} .
docker-push: docker-build
	docker tag ${DOCKER_IMAGE} steden/${DOCKER_IMAGE}:${VERSION}
	docker push steden/${DOCKER_IMAGE}:${VERSION}
