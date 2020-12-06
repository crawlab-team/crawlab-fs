FROM ubuntu:latest

RUN chmod 777 /tmp \
    && apt-get update \
    && apt-get install -y curl wget

RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.13/linux_amd64.tar.gz \
    && tar -zxf linux_amd64.tar.gz \
    && mv weed /usr/local/bin

