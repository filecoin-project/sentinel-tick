FROM golang:1.15-buster AS builder
MAINTAINER Hector Sanjuan <hector@protocol.ai>

ENV GOPATH      /go
ENV SRC_PATH    $GOPATH/src/github.com/filecoin-project/sentinel-tick
ENV GO111MODULE on
ENV GOPROXY     https://proxy.golang.org

RUN apt-get update && apt-get install -y ca-certificates

COPY go.* $SRC_PATH/
WORKDIR $SRC_PATH
RUN go mod download

COPY . $SRC_PATH
RUN go install -trimpath -mod=readonly -ldflags "-X main.tag=$(git describe)"

#-------------------------------------------------------------------

#------------------------------------------------------
FROM busybox:1-glibc
MAINTAINER Hector Sanjuan <hector@protocol.ai>

ENV GOPATH                 /go
ENV SRC_PATH    $GOPATH/src/github.com/filecoin-project/sentinel-tick

COPY --from=builder $GOPATH/bin/sentinel-tick /usr/local/bin/sentinel-tick
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/usr/local/bin/sentinel-tick"]

# Defaults for ipfs-cluster-service go here
CMD [""]
