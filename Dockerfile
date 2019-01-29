FROM golang:alpine3.8 AS builder
RUN apk add --no-cache make git ca-certificates
WORKDIR /go/src/github.com/shibayu36/notify-issues-to-slack/
COPY . .
RUN make build

FROM alpine:3.8
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/shibayu36/notify-issues-to-slack/notify-issues-to-slack /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/notify-issues-to-slack"]
