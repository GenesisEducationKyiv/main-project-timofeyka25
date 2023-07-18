# Build stage
FROM golang:1.20-alpine as build
LABEL maintainer="rate-service"
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main src/cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=build /app/main /main
COPY .env .env
COPY src/storage/csv src/storage/csv
CMD ["./main"]