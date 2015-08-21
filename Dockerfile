FROM centurylink/ca-certs

ADD ./bin/* /notification/
ADd config.yml /notification/

WORKDIR /notification

EXPOSE 3030

CMD ["./main","--config","config.yml"]

