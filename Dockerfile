FROM golang:1.10.3 AS builder

RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    echo '31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315  dep' | sha256sum -c - && \
    mv dep /usr/bin  && chmod 755 /usr/bin/dep

WORKDIR /go/src/github.com/nais/vault-kubernetes-secrets
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .
RUN \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -a -installsuffix cgo -o vks .

FROM scratch
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/nais/vault-kubernetes-secrets /
CMD ["/vks"]
