# Build stage
FROM golang:1.20-alpine as build
LABEL maintainer="logger-service"
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main src/cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=build /app/main /main
COPY .env .env
CMD ["./main"]