FROM golang:1.17-buster AS src

WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./server ./cmd/go8/main.go;


FROM scratch
LABEL MAINTAINER User <author@example.com>

WORKDIR /home/appuser/app/

COPY --from=src /go/src/app/server /home/appuser/app/server
COPY --from=src /go/src/app/.env /home/appuser/app/.env

EXPOSE 3080

ENTRYPOINT ["/home/appuser/app/server"]
