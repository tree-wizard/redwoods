ARG GO_VERSION=1.17
ARG ALPINE_VERSION=3.14

#Get official alpine image
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

#Install packages and Go language server bash   Make
RUN apk add -q --update --progress --no-cache sudo make  bash 
#GIT & nano & alpine-sdk 
RUN apk add curl git nano alpine-sdk 

WORKDIR /
Copy ./redwoods /go/bin/redwoods
RUN chmod +x /go/bin/redwoods
ENV SHELL=/bin/bash
CMD ["/bin/bash"]