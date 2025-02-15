# docker build -f e2e/server.Dockerfile .
FROM golang:1.24 AS src_server

# Copy dependencies first to take advantage of Docker caching
WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0

RUN go build -ldflags="-s" -o ./server ./cmd/go8/main.go;

FROM gcr.io/distroless/static-debian11:nonroot


LABEL com.gmhafiz.maintainers="User <author@example.com>"

WORKDIR /usr/local/bin

COPY --from=src_server /go/src/app/server .
COPY e2e/.env .env

#RUN apt update && apt install -y curl # curl is for healthcheck

EXPOSE 3090

CMD /usr/local/bin/server
