# Usage:
#   docker build -t chaitin/safeline-detector:t1k-ci-1.6.0 .
#   docker push chaitin/safeline-detector:t1k-ci-1.6.0
FROM chaitin/safeline-detector:1.6.0

RUN sed -i "s/^# bind_addr/bind_addr/; s/^# listen_port/listen_port/; s/^bind_addr: unix/# bind_addr: unix/;" /detector/snserver.yml
