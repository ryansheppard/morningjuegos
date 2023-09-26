all: test build push update

export DOCKER_IMAGE=registry.digitalocean.com/ryansheppard/morningjuegos
export DOCKER_TAG=$(shell git rev-parse --short HEAD)
export DOCKER_IMAGE_TAG=$(DOCKER_IMAGE):$(DOCKER_TAG)

test:
	go test ./... -v -cover
build: test
	docker build . -t ${DOCKER_IMAGE_TAG}
push: build
	docker push ${DOCKER_IMAGE_TAG}
update: build
	 yq e -i '.spec.template.spec.containers[0].image = "${DOCKER_IMAGE_TAG}"' deployments/manifests/deployment.yaml
