FROM docker.io/library/golang:1.13-alpine
RUN apk add --no-cache make upx
WORKDIR /build
COPY * ./
RUN make build \
 && make pack \
 && mv /build/longshore /longshore

FROM scratch
COPY --from=0 /longshore /longshore
ENTRYPOINT ["/longshore"]
