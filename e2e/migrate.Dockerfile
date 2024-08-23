FROM golang:1.22 AS src_migrate

# Copy dependencies first to take advantage of Docker caching
WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -o ./migrate ./cmd/migrate/main.go;

FROM gcr.io/distroless/static-debian11

LABEL com.gmhafiz.maintainers="User <author@example.com>"

WORKDIR /usr/local/bin

COPY --from=src_migrate /go/src/app/migrate .
COPY e2e/.env .env

CMD /usr/local/bin/migrate
