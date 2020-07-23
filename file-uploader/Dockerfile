FROM golang:1.14.4-alpine3.12

RUN groupadd -r sample && useradd --no-log-init -r -g sample sample

RUN mkdir "/src"
WORKDIR "/src"

COPY main main

CMD "./main"