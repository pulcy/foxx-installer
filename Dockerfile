FROM alpine:3.3

ADD ./foxx-installer /installer/

ENTRYPOINT ["/installer/foxx-installer"]
