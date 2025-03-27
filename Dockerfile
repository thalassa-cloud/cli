FROM gcr.io/distroless/static:nonroot

COPY tcloud /usr/local/bin/tcloud
ENTRYPOINT ["/usr/local/bin/tcloud"]
