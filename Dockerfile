FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o rate-limit ./cmd/server

FROM scratch
COPY --from=builder /app/rate-limit .
# Add the .env file but it's not the right way
COPY --from=builder /app/.env .env
CMD ["./rate-limit"]
