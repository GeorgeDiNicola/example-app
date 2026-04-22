DOCKER_IMAGE := my-go-app
PORT := 8080

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run -p $(PORT):$(PORT) $(DOCKER_IMAGE)

.PHONY: up
up: \
	docker-build \
	docker-run