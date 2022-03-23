FROM golang:1.17.5-alpine3.15 as builder
ADD . /go/host_exporter/
WORKDIR /go/host_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/host_exporter

FROM alpine:latest
ENV CONFIG_FILE "/config.yml"
ENV CMD_FLAGS ""

FROM centos

RUN dnf -y update \
    && dnf -y install freeipmi wget \
    && dnf clean all && rm -rf /var/cache/dnf 

WORKDIR /app
COPY --from=builder /go/bin/host_exporter .
CMD ./host_exporter --config.path $CONFIG_FILE $CMD_FLAGS
