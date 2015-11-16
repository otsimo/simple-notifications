FROM centurylink/ca-certs
MAINTAINER Sercan Degirmenci <sercan@otsimo.com>

ADD ./bin/notification-linux-amd64 /notification/simple-notifications
#ADD config.yml /notification/

WORKDIR /notification

#EXPOSE 18844

#CMD ["./simple-notifications","--config","config.yml"]
