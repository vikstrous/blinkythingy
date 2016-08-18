FROM golang:1.6.2-alpine

RUN apk add -U gcc linux-headers git libc-dev ca-certificates && \
        go get github.com/vikstrous/blinkythingy && \
        go install github.com/vikstrous/blinkythingy/...

ENTRYPOINT /go/bin/blinkythingy
