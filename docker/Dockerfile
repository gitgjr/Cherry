FROM golang:latest 

RUN apt update \
    && wget https://dist.ipfs.tech/kubo/v0.25.0/kubo_v0.25.0_linux-amd64.tar.gz \
    && tar xvfz kubo_v0.25.0_linux-amd64.tar.gz \
    && ./kubo/install.sh

RUN rm kubo_v0.25.0_linux-amd64.tar.gz\
    && rm -R ./kubo

RUN IPFS_PATH=~/.ipfs ipfs init 