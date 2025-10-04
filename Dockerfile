FROM debian:bookworm-slim

WORKDIR /app

# copy the pre-built binary from goreleaser's dist directory
# goreleaser expect an already build dist, so no need to build it again
COPY test .

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

CMD ["./go-actions-test"]