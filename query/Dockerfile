# Builder
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o query-service ./cmd/query-service

RUN chmod +x /app/query-service

# Runner
FROM scratch

COPY --from=builder /app/query-service .

CMD [ "./query-service" ]