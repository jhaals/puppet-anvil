FROM golang:onbuild
MAINTAINER Johan Haals <johan@haals.se>

ENV MODULEPATH /modules
ENV PORT 8080
EXPOSE 8080
