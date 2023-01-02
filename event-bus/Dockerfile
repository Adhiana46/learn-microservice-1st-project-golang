# Builder
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o event-bus ./cmd/event-bus

RUN chmod +x /app/event-bus

# Runner
FROM scratch

COPY --from=builder /app/event-bus .

CMD [ "./event-bus" ]