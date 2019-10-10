FROM golang:1.13 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine
WORKDIR /app

COPY --from=build-env /app /app

ENTRYPOINT ["./app"]