FROM alpine:latest

MAINTAINER Eric Mieth <mam09bog@studserv.uni-leipzig.de>

COPY mensabot /mensabot
COPY config.json /config.json

ENTRYPOINT /mensabot
