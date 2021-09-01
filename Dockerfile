FROM alpine:3.14
ENTRYPOINT ["/usr/local/bin/jsonnet-tool"]
COPY jsonnet-tool /usr/local/bin
