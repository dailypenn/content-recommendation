FROM golang:1.20-bookworm AS builder
WORKDIR /go/app
COPY . .
RUN go build -buildvcs=false

FROM redis/redis-stack-server:latest AS redis-stack
FROM redis:7-bookworm AS redis
COPY . .
COPY --from=redis-stack /opt/redis-stack/lib/redisearch.so /opt/redis-stack/lib/redisearch.so
COPY --from=builder /go/app/content-recommendation .
RUN apt-get update && apt-get install -y python3 python3-pip curl procps
RUN rm /usr/lib/python3.11/EXTERNALLY-MANAGED
RUN pip3 install pipenv
RUN pipenv install --system --deploy

ENV SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.2.25/supercronic-linux-amd64 \
    SUPERCRONIC=supercronic-linux-amd64 \
    SUPERCRONIC_SHA1SUM=642f4f5a2b67f3400b5ea71ff24f18c0a7d77d49

RUN curl -fsSLO "$SUPERCRONIC_URL" \
 && echo "${SUPERCRONIC_SHA1SUM}  ${SUPERCRONIC}" | sha1sum -c - \
 && chmod +x "$SUPERCRONIC" \
 && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
 && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic

CMD fallocate -l 512M /swapfile && \
    chmod 0600 /swapfile && \
    mkswap /swapfile && \
    echo 10 > /proc/sys/vm/swappiness && \
    swapon /swapfile && \
    echo 1 > /proc/sys/vm/overcommit_memory && \
    python3 server.py
