FROM alpine:3.14
ENTRYPOINT ["/jsonnet-tool"]
COPY jsonnet-tool /
