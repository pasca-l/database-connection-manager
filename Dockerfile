FROM ubuntu:latest

WORKDIR /home/local/

RUN apt-get update && apt-get install -y postgresql-client-16 mysql-client-8.0

COPY bin/* /usr/local/bin/
