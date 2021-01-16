FROM golang:1.15-alpine AS src

#RUN set -ex; \
#    apk update; \
#    apk --no-cache add ca-certificates git

WORKDIR /go/src/app/
COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./server ./cmd/go8/main.go;


#FROM debian:buster-slim
FROM alpine
LABEL MAINTAINER Hafiz <author@example.com>

#RUN apt update
#RUN apt install net-tools htop

# Add new user 'appuser'. App should be run without root privileges as a security measure
RUN adduser --home "/home/appuser" --disabled-password appuser --gecos "appuser,-,-,-"
USER appuser
RUN mkdir -p /home/appuser/app
WORKDIR /home/appuser/app/


#WORKDIR /opt/
COPY --from=src /go/src/app/server .
COPY .env .env

EXPOSE 3080

# Run Go Binary
#CMD /opt/server
CMD /home/appuser/app/server
