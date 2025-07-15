FROM golang:1.24.4-alpine as builded

WORKDIR /app

ARG ENABLE_RACE="false"

RUN if [ "$ENABLE_RACE" = "true" ]; then apk add --no-cache build-base; fi

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN if [ "$ENABLE_RACE" = "true" ]; then \
    export CGO_ENABLED=1 && go build -race -o main ./cmd/app/main.go; \
    else \
    export CGO_ENABLED=0 && go build -o main ./cmd/app/main.go; \
    fi

FROM alpine

WORKDIR /app

COPY --from=builded /app/main /

ARG API_PORT

EXPOSE ${API_PORT}

ENTRYPOINT /main
