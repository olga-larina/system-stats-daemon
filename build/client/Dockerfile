# Собираем в гошке
FROM golang:1.22 AS build

ENV BIN_FILE=/opt/stats/client
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
RUN go mod tidy

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/client/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="system-stats-client"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE="/opt/stats/client"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/stats/client.yaml

CMD ${BIN_FILE} -config ${CONFIG_FILE}
