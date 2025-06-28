FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . .

RUN go build -o ./app ./cmd/main.go

FROM alpine

WORKDIR /

COPY --from=builder /app/app /app

CMD ["/app"]