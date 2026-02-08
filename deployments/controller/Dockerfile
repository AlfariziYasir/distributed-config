FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o controller cmd/controller/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir -p /app/data
COPY --from=builder /app/controller .
RUN chmod -R 777 /app/data
EXPOSE 8080
CMD ["./controller"]