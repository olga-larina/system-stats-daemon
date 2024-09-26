# Собираем в гошке
FROM golang:1.22 AS build

ENV BIN_FILE=/opt/stats/daemon
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
RUN go mod tidy

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API).
# Используем build-tag linux.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -tags linux \
        -o ${BIN_FILE} cmd/daemon/*

# На выходе - образ с ubuntu
FROM ubuntu:24.04

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="system-stats-daemon"
LABEL MAINTAINERS="student@otus.ru"

# установка таймзоны (для iostat)
ENV TZ=Europe/London
RUN ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone

RUN apt-get update && \
    apt-get install -y --no-install-recommends tzdata sysstat wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# скачиваем healthcheck grpc
ENV GRPC_HEALTH_PROBE="https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.11/grpc_health_probe-linux-amd64"
ENV GRPC_HEALTH_PROBE_FILE="/usr/local/bin/grpc_health_probe"
RUN wget -q --no-check-certificate -o ${GRPC_HEALTH_PROBE_FILE} ${GRPC_HEALTH_PROBE}
RUN chmod +x ${GRPC_HEALTH_PROBE_FILE}

ENV BIN_FILE="/opt/stats/daemon"

# Наш сервис
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/stats/daemon.yaml

CMD ${BIN_FILE} -config ${CONFIG_FILE}
