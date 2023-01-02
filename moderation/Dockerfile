# Builder
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o moderation-service ./cmd/moderation-service

RUN chmod +x /app/moderation-service

# Runner
FROM scratch

COPY --from=builder /app/moderation-service .

CMD [ "./moderation-service" ]