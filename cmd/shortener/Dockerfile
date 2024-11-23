FROM golang:1.22.7 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /shortener ./cmd/shortener/main.go

FROM alpine AS server
WORKDIR /
COPY --from=builder /shortener /shortener
EXPOSE 8080
ENTRYPOINT ["/shortener"]