FROM alpine:3.20
ENTRYPOINT ["/usr/local/bin/jsonnet-tool"]
COPY jsonnet-tool /usr/local/bin
