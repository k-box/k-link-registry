FROM node:8
WORKDIR /app
ADD ui .
RUN \
    npm install && \
    npm run build

FROM golang:1.10
WORKDIR /go/src/github.com/k-box/k-link-registry
ADD . .
COPY --from=0 /app/dist ui/dist
RUN \
    go get -tags="dev" -v github.com/k-box/k-link-registry/klinkregistry && \
    go get github.com/shurcooL/vfsgen/cmd/vfsgendev && \
    go generate github.com/k-box/k-link-registry/assets && \
    cd klinkregistry && go build -tags="netgo" -o klinkregistry

FROM scratch
EXPOSE 80
COPY --from=1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=1 /go/src/github.com/k-box/k-link-registry/klinkregistry/klinkregistry /klinkregistry
ENTRYPOINT ["/klinkregistry"]
CMD ["server"]
