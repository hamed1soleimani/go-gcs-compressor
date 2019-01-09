FROM golang:1.10-alpine3.8 AS build

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml $GOPATH/src/github.com/hamed1soleimani/go-gcs-compressor/
WORKDIR $GOPATH/src/github.com/hamed1soleimani/go-gcs-compressor
RUN dep ensure -vendor-only

COPY main.go  $GOPATH/src/github.com/hamed1soleimani/go-gcs-compressor/
RUN CGO_ENABLED=0 GOOS=linux go build -o $GOPATH/bin/go-gcs-compressor -a -ldflags '-extldflags "-static"' main.go
RUN apk add -U --no-cache ca-certificates

FROM scratch
ENV GOPATH /go

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build $GOPATH/bin $GOPATH/bin
ENTRYPOINT ["/go/bin/go-gcs-compressor"]
