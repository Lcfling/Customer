# Dockerfile
FROM centos:centos8
MAINTAINER Lcfling <aylui2009@163.com>
WORKDIR /www/go/
ADD ./copy/ /www/go/
#RUN go build .
EXPOSE 23001 13892 8077 8843
EXPOSE 61005/udp
#ENTRYPOINT ["./bin/imr","-l","-log_dir=./logs/imr","-c","imr.cfg"]
ENTRYPOINT ./main -log_dir=./logs/game im.cfg