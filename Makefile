IMAGE_NAME ?= example-app
PORT ?= 8080


.PHONY: test
test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	docker buildx build -t $(IMAGE_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME)

.PHONY: up
up: \
	docker-build \
	docker-run

.PHONY: docker-test
docker-test:
	# Run the container in the background
	docker run -d --name health-check -p $(PORT):$(PORT) $(IMAGE_NAME)
	sleep 3
	# force an error if curl fails the healthcheck
	curl --fail http://localhost:$(PORT)/health || (docker logs health-check && docker stop health-check && docker rm health-check && exit 1)
	docker stop health-check
	docker rm health-check