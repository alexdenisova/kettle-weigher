ARG VERSION="3.16"
FROM alpine:${VERSION} AS final

LABEL org.opencontainers.image.source = "https://github.com/alexdenisova/kettle-weigher"

ARG USER="user"
ARG HOME="/app"
ARG UID="1000"
ARG GID="1000"
RUN \
  addgroup --gid "${GID}" "${USER}" \
  ; adduser --disabled-password --gecos "" --home "${HOME}" --ingroup "${USER}" --uid "${UID}" "${USER}"

RUN \
  apk update \
  && apk --no-cache add busybox-extras curl jq vim \
  && rm -rf /var/cache/apk/* /tmp/*

COPY dist/bin/* /bin/
RUN chmod +x /bin/*

USER ${USER}
EXPOSE 8080/tcp
WORKDIR ${HOME}
ENTRYPOINT ["kettle-weigher"]
