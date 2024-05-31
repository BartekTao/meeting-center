FROM golang:1.22-alpine3.19 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY tmp.env .env 
RUN go build --o /app ./cmd/api/server.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app .
COPY --from=builder /build/.env .
CMD [ "./app" ]