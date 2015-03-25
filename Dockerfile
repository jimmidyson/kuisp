FROM gliderlabs/alpine:3.1
MAINTAINER Jimmi Dyson <jimmidyson@gmail.com>

ENV VERSION 0.1

RUN apk-install ca-certificates curl tar && \
    curl -L https://github.com/jimmidyson/kuisp/releases/download/v0.1/kuisp-0.1-linux-amd64.tar.gz | \
      tar xzv && \
    apk del curl tar

ENTRYPOINT ["/kuisp"]
