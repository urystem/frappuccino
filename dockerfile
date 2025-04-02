FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN go build -o frappuccino cmd/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /app/frappuccino .

# EXPOSE 8080

ENTRYPOINT [ "./frappuccino" ]
