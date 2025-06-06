# -------- Dockerfile --------
FROM golang:1.23-alpine as builder

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN GOOS=linux go build -o bankingapp ./cmd


# -- runtime image --
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bankingapp .
ADD configs/app_config.json configs/app_config.json

EXPOSE 8080

CMD ["./bankingapp"]