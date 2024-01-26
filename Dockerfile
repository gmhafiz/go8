FROM golang:1.21 AS src

WORKDIR /go/src/app/

# Copy dependencies first to take advantage of Docker caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Insert version using git tag and latest commit hash
# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=$(git describe --abbrev=0 --tags)-$(git rev-list -1 HEAD) -w -s" -o ./server ./cmd/go8/main.go;

FROM gcr.io/distroless/static-debian12:nonroot

LABEL com.example.maintainers="User <author@example.com>"

COPY --from=src /go/src/app/server /usr/bin/local/server

EXPOSE 3080

ENTRYPOINT ["/usr/bin/local/server"]
