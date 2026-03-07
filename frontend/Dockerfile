# frontend/Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY script.js style.css index.html env.js send_mes.png ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/app .
EXPOSE 8082
CMD ["./app"]