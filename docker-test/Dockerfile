FROM golang:1.15-alpine AS src

#RUN set -ex; \
#    apk update; \
#    apk --no-cache add ca-certificates git

WORKDIR /go/src/app/
COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./server ./cmd/go8/main.go;

RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./e2e ./cmd/e2e/main.go;

FROM debian:buster-slim
#FROM alpine
LABEL MAINTAINER Hafiz <author@example.com>

RUN apt update
RUN apt install -y net-tools htop postgresql-client nano curl

# Add new user 'appuser'. App should be run without root privileges as a security measure
#RUN adduser --home "/home/appuser" --disabled-password appuser --gecos "appuser,-,-,-"
#USER appuser
RUN useradd -m -d /home/appuser appuser
RUN id -u appuser | awk '{print $1+1}' | xargs -I{} usermod -u {} appuser
RUN mkdir -p /home/appuser/app
WORKDIR /home/appuser/app/


#WORKDIR /opt/
RUN mkdir -p database/migrations
COPY --from=src /go/src/app/server .
COPY --from=src /go/src/app/e2e .
COPY --from=src /go/src/app/database/migrations/ database/migrations/
COPY ./docker-test/.env .env
RUN id -u appuser | xargs -I{} chown -R {}:{} /home/appuser/app

EXPOSE 3081

# Run Go Binary
CMD cd /home/appuser/app/ && /home/appuser/app/e2e
#CMD /home/appuser/app/e2e
#CMD echo "hello"
