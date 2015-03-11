##Dockerizing a go web app
FROM golang

MAINTAINER Kristina Hager <kristina.hager@gmail.com>

#don't include git history in this image
RUN echo ".git" > .dockerignore

#mgo for talking to mongodb
RUN go get gopkg.in/mgo.v2

COPY . /go/src/test-docker-go-mongodb
RUN go install test-docker-go-mongodb

ENTRYPOINT /go/bin/test-docker-go-mongodb
EXPOSE 8080
