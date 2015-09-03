FROM centurylink/ca-certs

ADD ./bin/* /notification/
ADD config.yml /notification/

WORKDIR /notification

EXPOSE 3030

CMD ["./main-linux-amd64","--config","config.yml"]
