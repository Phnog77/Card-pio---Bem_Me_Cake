FROM golang:1.25.2 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

RUN CGO_ENABLED=0 go build -o server ./app

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .

COPY templates ./templates
COPY static ./static

EXPOSE 443

CMD ["./server"]
