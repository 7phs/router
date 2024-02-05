FROM golang:1.21-alpine as builder
WORKDIR /app

COPY go.* ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go mod download
RUN go build -v -o ./bin/ ./cmd/*

FROM alpine

COPY --from=builder /app/bin/* /app/

EXPOSE 8080
CMD ["/app/router"]