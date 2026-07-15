# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/laundry ./cmd/server/

# Stage 2: Run
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/laundry .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config.yaml .

RUN chown -R app:app /app

USER app

EXPOSE 3000

CMD ["./laundry"]
