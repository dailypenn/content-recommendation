FROM golang:1.20-bookworm AS builder
WORKDIR /go/app
COPY . .
RUN go build -buildvcs=false

FROM redis/redis-stack-server:latest AS redis-stack
FROM redis:7-bookworm AS redis
COPY . .
COPY --from=redis-stack /opt/redis-stack/lib/redisearch.so /opt/redis-stack/lib/redisearch.so
COPY --from=builder /go/app/content-recommendation .
RUN apt-get update && apt-get install -y python3 python3-pip curl procps cron vim
ADD crontab /etc/cron.d/cronjobs 
RUN chmod 0644 /etc/cron.d/cronjobs && \ 
    crontab /etc/cron.d/cronjobs && \
    touch /var/log/cron.log
RUN rm /usr/lib/python3.11/EXTERNALLY-MANAGED
RUN pip3 install pipenv
RUN pipenv install --system --deploy

CMD fallocate -l 512M /swapfile && \
    chmod 0600 /swapfile && \
    mkswap /swapfile && \
    echo 10 > /proc/sys/vm/swappiness && \
    swapon /swapfile && \
    echo 1 > /proc/sys/vm/overcommit_memory && \
    python3 server.py --refresh
