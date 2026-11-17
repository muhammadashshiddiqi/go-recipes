FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /recipe-api .

FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=builder /recipe-api /recipe-api
EXPOSE 8080
ENTRYPOINT ["/recipe-api"]