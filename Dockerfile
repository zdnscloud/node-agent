FROM golang:1.12.5-alpine3.9 AS build

RUN mkdir -p /go/src/github.com/zdnscloud/node-agent
COPY . /go/src/github.com/zdnscloud/node-agent

WORKDIR /go/src/github.com/zdnscloud/node-agent
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/zdnscloud/node-agent/node-agent

FROM alpine:3.9.4

LABEL maintainers="Zdns Authors"
LABEL description="Node Agent"
RUN apk update && apk add util-linux udev --no-cache
COPY --from=build /go/src/github.com/zdnscloud/node-agent/node-agent /node-agent
ENTRYPOINT ["/bin/sh"]
