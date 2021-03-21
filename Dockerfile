FROM golang:1.16.2-buster

RUN apt-get update &&\
    apt-get install -y python3 python3-yaml &&\
    apt-get clean &&\
    rm /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

ARG CURRENT_TAG
ENV CURRENT_TAG=${CURRENT_TAG}
RUN python3 scripts/build.py
