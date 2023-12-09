FROM alpine:3.19
ENTRYPOINT ["/usr/local/bin/jsonnet-tool"]
COPY jsonnet-tool /usr/local/bin
