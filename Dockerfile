FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/

FROM alpine:latest

ARG USER=appuser

RUN adduser -D -g '' $USER

WORKDIR /app

COPY --from=builder /app/main .

RUN chown -R $USER /app

USER $USER

CMD ["./main"]