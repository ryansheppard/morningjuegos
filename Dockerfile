FROM golang:1.21 AS builder

WORKDIR /src
COPY . .

RUN apt-get update && apt-get install -y git
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch

COPY --from=builder /go/bin/morningjuegos /morningjuegos
COPY --from=golang:1.21 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/morningjuegos"]
