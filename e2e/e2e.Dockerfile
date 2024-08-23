# docker build -f e2e/e2e.Dockerfile -t go8/e2e .
# docker run -it go8/e2e
FROM golang:1.22 AS src_e2e

WORKDIR /go/src/app/

# Copy dependencies first to take advantage of Docker caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -o ./end_to_end ./e2e/main.go;


FROM gcr.io/distroless/static-debian11

LABEL com.example.maintainers="User <author@example.com>"

WORKDIR /usr/local/bin

COPY --from=src_e2e /go/src/app/end_to_end /usr/local/bin/end_to_end
COPY e2e/.env .env


ENTRYPOINT ["/usr/local/bin/end_to_end"]
