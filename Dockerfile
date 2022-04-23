FROM golang AS builder
WORKDIR $GOPATH/src/github.com/jeffre/downloader
COPY . .
RUN update-ca-certificates \
    && go test ./... \
    && CGO_ENABLED=0 go build -o /tmp/downloader cmd/downloader/downloader.go

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp/downloader /
ENTRYPOINT [ "/downloader"]
CMD ["-h"]