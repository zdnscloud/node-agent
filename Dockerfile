FROM golang:1.13.7-alpine3.11 AS build
ENV GOPROXY=https://goproxy.cn

RUN mkdir -p /go/src/github.com/zdnscloud/node-agent
COPY . /go/src/github.com/zdnscloud/node-agent

WORKDIR /go/src/github.com/zdnscloud/node-agent
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/zdnscloud/node-agent/node-agent

FROM ubuntu:14.04
LABEL maintainers="Zdns Authors"
LABEL description="Node Agent"

RUN apt-get update && apt-get install util-linux udev open-iscsi lvm2 lsscsi multipath-tools -y
COPY --from=build /go/src/github.com/zdnscloud/node-agent/node-agent /node-agent
ENTRYPOINT ["/bin/sh"]
