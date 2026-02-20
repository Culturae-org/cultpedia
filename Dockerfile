FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION="1.0.0"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'cultpedia/internal/utils.Version=${VERSION}'" -a -installsuffix cgo -o cultpedia cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/cultpedia .

COPY --from=builder /app/datasets ./datasets

EXPOSE 8080

CMD ["./cultpedia", "api"]