FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /collector ./collector/cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /collector .
EXPOSE 50051
CMD ["./collector"]
