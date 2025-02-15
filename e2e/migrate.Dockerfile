FROM golang:1.24 AS src_migrate

# Copy dependencies first to take advantage of Docker caching
WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0

RUN go build -ldflags="-s" -o ./migrate ./cmd/migrate/main.go;

FROM gcr.io/distroless/static-debian12:nonroot

LABEL com.gmhafiz.maintainers="User <author@example.com>"

WORKDIR /usr/local/bin

COPY --from=src_migrate /go/src/app/migrate .
COPY e2e/.env .env

CMD /usr/local/bin/migrate
