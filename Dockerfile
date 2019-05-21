FROM golang:alpine AS build
  
RUN mkdir -p /go/src/github.com/zdnscloud/node-agent
COPY . /go/src/github.com/zdnscloud/node-agent

WORKDIR /go/src/github.com/zdnscloud/node-agent
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/src/github.com/zdnscloud/node-agent/node-agent

FROM alpine

LABEL maintainers="Zdns Authors"
LABEL description="Node Agent"
COPY --from=build /go/src/github.com/zdnscloud/node-agent/node-agent /node-agent
ENTRYPOINT ["/bin/sh"]
