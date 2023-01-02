# Builder
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o posts-service ./cmd/posts-service

RUN chmod +x /app/posts-service

# Runner
FROM scratch

COPY --from=builder /app/posts-service .

CMD [ "./posts-service" ]