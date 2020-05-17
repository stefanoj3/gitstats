FROM golang:1.14.3-alpine3.11 as builder

RUN adduser -D -g '' gitstatsuser
RUN apk add --update make git ca-certificates

RUN mkdir -p $GOPATH/src/github.com/stefanoj3/gitstats
ADD . $GOPATH/src/github.com/stefanoj3/gitstats
WORKDIR $GOPATH/src/github.com/stefanoj3/gitstats

RUN make build

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/stefanoj3/gitstats/dist/gitstats /bin/gitstats
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER gitstatsuser
CMD ["gitstats"]