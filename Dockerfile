FROM golang:1.13 AS build

WORKDIR /go/src/github.com
ENV GOPATH=/go PATH=$PATH:/go/bin
RUN git clone --single-branch --branch master https://github.com/protosio/protos.git
WORKDIR /go/src/github.com/protos
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o bin/protos cmd/protos/protos.go
RUN mkdir /root/tmp


FROM alpine:3.10.3
WORKDIR /
COPY --from=build /go/src/github.com/protos/bin/protos /usr/local/bin/protos
RUN chmod +x /usr/local/bin/protos
RUN mkdir /var/protos && mkdir /var/protos-containerd
COPY protos.yaml /etc/protos.yaml

ENTRYPOINT ["/usr/local/bin/protos", "--loglevel", "debug", "--config", "/etc/protos.yaml", "daemon"]
