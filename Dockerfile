FROM alpine:3.18
ENTRYPOINT ["/usr/local/bin/jsonnet-tool"]
COPY jsonnet-tool /usr/local/bin
