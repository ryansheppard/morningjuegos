all: test build push

test:
	go test ./... -v -cover
build: test
	docker build . -t ${DOCKER_IMAGE}
push: build
	docker push ${DOCKER_IMAGE}
