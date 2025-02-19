ARG ALPINE_VERSION="3.16"
FROM alpine:${ALPINE_VERSION}

RUN \
  apk update \
  && apk --no-cache add vim \
  && rm -rf /var/cache/apk/* /tmp/*

ARG GOLANG_VERSION="1.23.1"
RUN \
  wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && rm -f go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH=/usr/local/go/bin:${PATH}

WORKDIR /workspace
COPY cmd/ ./cmd
COPY go.mod ./
RUN \
  cd cmd \
  && go build -o kettle-weigher . \
  && mv kettle-weigher /usr/local/bin/kettle-weigher \
  && chmod 0755 /usr/local/bin/kettle-weigher

ENTRYPOINT [""]
