FROM ubuntu
MAINTAINER Johan Haals <johan@haals.se>

RUN apt-get update
RUN apt-get install -y golang

ADD . /source
RUN cd /source && go build -o /puppet-anvil

ENV MODULEPATH /modules
ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["/puppet-anvil"]
