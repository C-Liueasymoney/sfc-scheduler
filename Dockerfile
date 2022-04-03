FROM debian:stretch-slim

WORKDIR /

COPY bin/sfc-scheduler /usr/local/bin

CMD ["sfc-scheduler"]
