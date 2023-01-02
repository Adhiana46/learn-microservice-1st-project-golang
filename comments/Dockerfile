# Builder
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o comments-service ./cmd/comments-service

RUN chmod +x /app/comments-service

# Runner
FROM scratch

COPY --from=builder /app/comments-service .

CMD [ "./comments-service" ]