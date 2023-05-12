## Build
FROM golang:1.20.2-alpine3.16 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/ ./cmd/app/

## Deploy
FROM alpine:3.14

WORKDIR /app

COPY --from=build /bin .

EXPOSE 8000