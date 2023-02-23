FROM alpine:3.17
ENTRYPOINT ["/usr/local/bin/jsonnet-tool"]
COPY jsonnet-tool /usr/local/bin
