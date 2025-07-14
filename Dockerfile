FROM golang:1.24.4-alpine as builded

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build cmd/app/main.go

FROM alpine

WORKDIR /app

COPY --from=builded /app/main /

ARG API_PORT

EXPOSE ${API_PORT}

ENTRYPOINT /main
