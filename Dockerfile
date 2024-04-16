FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o webcoinApp ./cmd/webcoin

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/webcoinApp /app/webcoinApp

CMD [ "/app/webcoinApp" ]