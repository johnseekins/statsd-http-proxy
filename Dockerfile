FROM debian:stable-slim

COPY statsd-http-proxy /usr/local/bin/
RUN apt-get update -qq \
  && apt-get install -y -qq --no-install-recommends curl \
  && apt-get clean \
  && apt-get autoremove -y \
  && chmod +x /usr/local/bin/statsd-http-proxy

# start service
EXPOSE 80
ENTRYPOINT ["/usr/local/bin/statsd-http-proxy", "--http-host="]
