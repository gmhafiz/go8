# docker build -f e2e/e2e.Dockerfile -t go8/e2e .
# docker run -it go8/e2e
FROM golang:1.20 AS src

# Copy dependencies first to take advantage of Docker caching
WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./end_to_end ./cmd/e2e/main.go;


FROM gcr.io/distroless/static-debian11

LABEL com.gmhafiz.maintainers="User <author@example.com>"

WORKDIR /usr/local/bin

COPY --from=src /go/src/app/end_to_end /usr/local/bin/end_to_end
COPY ./e2e/.env .env


# Run Go Binary
ENTRYPOINT ["/usr/local/bin/end_to_end"]
