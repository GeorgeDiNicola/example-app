IMAGE_NAME ?= example-app
PORT ?= 8080
VERSION_TAG ?= v1.0.0


.PHONY: test
test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	docker buildx build \
		-t $(IMAGE_NAME):latest \
		-t $(IMAGE_NAME):$(VERSION_TAG) \
		.

.PHONY: docker-run
docker-run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME)

.PHONY: up
up: \
	docker-build \
	docker-run

.PHONY: docker-test
docker-test:
	# prevent docker container name collisions by using the shell PID
	@container_name=health-check-$$$$; \
	docker run -d --name $$container_name -p $(PORT):$(PORT) $(IMAGE_NAME):$(VERSION_TAG); \
	trap 'docker logs $$container_name 2>/dev/null || true; docker rm -f $$container_name 2>/dev/null || true' 0 1 2 3 15; \
	for i in 1 2 3; do \
		if curl --fail --silent http://localhost:$(PORT)/health > /dev/null; then \
			exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "Smoke test failed"; \
	exit 1
