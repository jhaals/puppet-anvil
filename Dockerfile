FROM benschw/litefs

ADD ./puppet-anvil /puppet-anvil

ENV MODULEPATH /modules
ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["/puppet-anvil"]
