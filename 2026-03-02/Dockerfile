FROM golang:1.24-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o /out/server ./cmd/server && \
    go build -o /out/bench ./cmd/bench && \
    go build -o /out/compare ./cmd/compare

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/server /usr/local/bin/server
COPY --from=build /out/bench /usr/local/bin/bench
COPY --from=build /out/compare /usr/local/bin/compare
ENTRYPOINT ["server"]
