FROM ubuntu
MAINTAINER Johan Haals <johan@haals.se>

RUN apt-get update
RUN apt-get install -y golang make git-core

ADD . /go/src/github.com/jhaals/puppet-anvil
RUN export GOPATH=/go && \
		cd /go/src/github.com/jhaals/puppet-anvil && \
		make deps build && \
		mv /go/src/github.com/jhaals/puppet-anvil/build/puppet-anvil /puppet-anvil && \
		rm -rf /go

ENV MODULEPATH /modules
ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["/puppet-anvil"]
