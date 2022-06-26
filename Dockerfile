FROM golang:1.18.2-alpine3.16 as builder

ENV GOPATH=/go

COPY . $GOPATH/src/github.com/stockwayup/http

WORKDIR $GOPATH/src/github.com/stockwayup/http

RUN go get -u -t github.com/tinylib/msgp && \
    go install github.com/tinylib/msgp && \
    go generate ./... && \
    go build -o /bin/stockwayup

FROM alpine:3.16

COPY --from=builder --chown=www-data /bin/stockwayup /bin/stockwayup

RUN chmod +x /bin/stockwayup

USER www-data

CMD ["/bin/stockwayup"]