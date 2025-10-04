FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main/main.go


FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main . 
EXPOSE 8080
CMD ["/app/main"]