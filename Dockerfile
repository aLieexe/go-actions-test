FROM debian:bookworm-slim

WORKDIR /app

# copy the pre-built binary from goreleaser's dist directory
COPY go-actions-test .

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

CMD ["./go-actions-test"]