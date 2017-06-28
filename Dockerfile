FROM golang:1.8
LABEL version="0.2" \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

RUN  go get github.com/golang/dep/cmd/dep && \
     go get github.com/mhausenblas/reshifter && \
     mkdir -p /app/ui

WORKDIR /go/src/github.com/mhausenblas/reshifter
RUN dep ensure && \
    go build . && \
    mv reshifter /app/reshifter
COPY ui/ /app/ui/

USER nobody
EXPOSE 8080
CMD ["/app/reshifter"]
