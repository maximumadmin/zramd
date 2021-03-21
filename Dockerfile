FROM golang:1.16.2-buster

RUN apt-get update &&\
    apt-get install -y python3 python3-yaml &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

RUN go mod download

ARG CURRENT_TAG
ENV CURRENT_TAG=${CURRENT_TAG}
ARG COMMIT_DATE
ENV COMMIT_DATE=${COMMIT_DATE}

RUN python3 scripts/build.py
