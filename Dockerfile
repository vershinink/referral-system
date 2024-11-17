FROM golang:1.23.3-alpine3.20 AS builder

ENV REF_CONFIG_PATH=./config/dev.yaml

WORKDIR /go/src/referral

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/referral/main.go

FROM alpine:3.20 AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /root

ENV REF_CONFIG_PATH=/root/config/dev.yaml

RUN mkdir -p /root/config

COPY --from=builder /go/src/referral/config ./config

COPY --from=builder /go/src/referral/main .

EXPOSE 8080

RUN chmod +x /root/main

ENTRYPOINT [ "/root/main" ]