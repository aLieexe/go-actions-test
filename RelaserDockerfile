FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

ARG BINARY_PATH
COPY ${BINARY_PATH} /app/test
CMD ["./test"]