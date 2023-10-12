#FROM alpine:latest
FROM 127.0.0.1:5000/alpine-docker:develop

MAINTAINER  jzy
WORKDIR /go/src/ugk-agent

COPY ugk-agent /go/src/ugk-agent/ugk-agent
COPY config/*.json /go/src/ugk-agent/config/

CMD /go/src/ugk-agent/ugk-agent ${GO_OPTS}
