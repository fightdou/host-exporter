FROM golang:1.17.5-alpine3.15 as builder
ADD . /go/host_exporter/
WORKDIR /go/host_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/host_exporter cmd/main.go

FROM douyali/centos8-5-storcli:latest

ENV CONFIG_FILE "/opt/config.yml"
ENV CMD_FLAGS ""

WORKDIR /opt
COPY --from=builder /go/bin/host_exporter .
CMD ./host_exporter --config.path $CONFIG_FILE $CMD_FLAGS
