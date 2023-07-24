.PHONY: docker-build-server
docker-build-server:
	docker build -t pow_ddos_server --file Dockerfile.server .

.PHONY: docker-build-client
docker-build-client:
	docker build -t pow_ddos_client --file Dockerfile.client .


.PHONY: docker-run-server
docker-run-server:
	docker network create pow_ddos
	docker run --rm -d \
		--network pow_ddos \
		--env-file env.server.example \
		--name pow_ddos_server \
		pow_ddos_server

.PHONY: docker-run-client
docker-run-client:
	docker run --rm -it \
		--network pow_ddos \
		--env-file env.client.example \
		--name pow_ddos_client \
		pow_ddos_client

.PHONY: docker-down
docker-down:
	docker rm -f pow_ddos_client
	docker rm -f pow_ddos_server
	docker network rm pow_ddos
