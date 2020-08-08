FROM golang:1.14 AS builder

COPY ${PWD} /app
WORKDIR /app

# Toggle CGO on your app requirement
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin cmd/go8/go8.go
# Use below if using vendor
# RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-extldflags "-static"' -o /app/appbin *.go

FROM debian:stable-slim
LABEL MAINTAINER Hafiz <author@example.com>

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# Add new user 'appuser'. App should be run without root privileges as a security measure
RUN adduser --home "/home/appuser" --disabled-password appuser --gecos "appuser,-,-,-"
USER appuser

#COPY --from=builder /app/internal/server/http/web /home/appuser/app/web
RUN mkdir -p /home/appuser/app
RUN mkdir -p /home/appuser/app/config
COPY --from=builder /app/appbin /home/appuser/app
COPY --from=builder /app/config/*.yml /home/appuser/app/config/

#ENV TEMPLATES_BASEPATH=/home/appuser/app/web/templates

WORKDIR /home/appuser/app

# Since running as a non-root user, port bindings < 1024 are not possible
# 8000 for HTTP; 8443 for HTTPS;
EXPOSE 3080

CMD ["./appbin"]