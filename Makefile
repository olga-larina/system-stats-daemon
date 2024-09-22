BIN_DAEMON := "./bin/daemon"
DOCKER_IMG_DAEMON :="system-stats-daemon:develop"

BIN_CLIENT := "./bin/client"
DOCKER_IMG_CLIENT:="system-stats-client:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# скомпилировать бинарные файлы сервиса (под darwin)
.PHONY: build_daemon
build_daemon:
	go build -v -o $(BIN_DAEMON) -ldflags "$(LDFLAGS)" ./cmd/daemon

.PHONY: build_client
build_client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

.PHONY: build
build: build_daemon build_client

# собрать и запустить сервисы с конфигами по умолчанию
.PHONY: run_daemon
run_daemon: build_daemon
	$(BIN_DAEMON) -config ./configs/daemon_local.yaml

.PHONY: run_client
run_client: build_client
	$(BIN_CLIENT) -config ./configs/client_local.yaml

# собрать все образы
.PHONY: build-img
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_DAEMON) \
		-f build/daemon/Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_CLIENT) \
		-f build/client/Dockerfile .

# запустить все образы
.PHONY: run-img
run-img: build-img
	docker run $(DOCKER_IMG_DAEMON) -d
	docker run $(DOCKER_IMG_CLIENT) -d

# запустить юнит-тесты
.PHONY: test
test:
	go test -race -count 100 ./internal/...

# линтер golangci-lint 
.PHONY: install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...

# сгенерировать прото
# необходимо сначала выполнить
# - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# - go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# - создать пакет для кода grpc-сервера с go:generate
.PHONY: generate
generate: 
	go generate ./...

# поднять сервисы
.PHONY: up
up:
	docker compose --env-file deployments/.env -f deployments/docker-compose.yaml -f deployments/docker-compose-client.yaml up -d --build

# потушить сервисы
.PHONY: down
down:
	docker compose --env-file deployments/.env -f deployments/docker-compose.yaml -f deployments/docker-compose-client.yaml down

# запустить интеграционные тесты
.PHONY: integration-tests
integration-tests:
	bash integration_tests/integration_tests.sh; \
    e=$$?; \
    exit $$e