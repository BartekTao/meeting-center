FROM golang:1.22-alpine3.19 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --o /app ./cmd/cronjob/reminder.go

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app .

CMD [ "./app" ]