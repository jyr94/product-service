
FROM golang:1.23-alpine AS test

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test ./... -v


FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/


FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=build /app/app .

EXPOSE 8181

CMD ["./app"]
