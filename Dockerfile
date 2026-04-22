# STAGE 1: Build the binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Optimize Docker layer caching so packages are not downloaded on every rebuild
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd/server/main.go

# STAGE 2: Run the app
FROM alpine:3.23

RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/main .

USER appuser

EXPOSE 8080

CMD ["./main"]