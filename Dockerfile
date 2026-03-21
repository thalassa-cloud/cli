FROM docker.io/library/alpine:3.23

ARG TARGETPLATFORM

RUN apk add --no-cache ca-certificates \
	&& addgroup -S thalassa && adduser -S -G thalassa -u 65532 -h /home/thalassa thalassa

COPY ${TARGETPLATFORM}/tcloud /usr/local/bin/tcloud
RUN chmod 755 /usr/local/bin/tcloud

USER thalassa
WORKDIR /home/thalassa

ENTRYPOINT ["/usr/local/bin/tcloud"]
