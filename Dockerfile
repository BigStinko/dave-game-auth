FROM golang:1.22-alpine AS builder
WORKDIR /app
# go.sum later
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main
# alpine image for final container
FROM alpine:latest
WORKDIR /APP
COPY --from=builder /app/main .
EXPOSE 8000
CMD ["./main"]
