FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY handlers ./handlers
COPY middleware ./middleware
COPY templates ./templates
COPY styles ./styles

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/toctou-demo .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/toctou-demo ./toctou-demo
COPY --from=builder /src/templates ./templates
COPY --from=builder /src/styles ./styles

EXPOSE 8080

CMD ["./toctou-demo"]
