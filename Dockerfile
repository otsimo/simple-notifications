FROM centurylink/ca-certs

ADD ./bin/notification-linux-amd64 /notification/
ADD config.yml /notification/

WORKDIR /notification

EXPOSE 18844

CMD ["./notification-linux-amd64","--config","config.yml"]
