FROM alpine:3.4
MAINTAINER Sercan Degirmenci <sercan@otsimo.com>

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

ADD notification-linux-amd64 /opt/otsimo/simple-notifications

EXPOSE 18844

CMD ["/opt/otsimo/simple-notifications"]
