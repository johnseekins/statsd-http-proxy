# StatsD HTTP Proxy

StatsD HTTP proxy with REST interface for using in browsers

[![Go Report Card](https://goreportcard.com/badge/github.com/GoMetric/statsd-http-proxy?2)](https://goreportcard.com/report/github.com/GoMetric/statsd-http-proxy)
[![Build Status](https://travis-ci.org/GoMetric/statsd-http-proxy.svg?branch=master)](https://travis-ci.org/GoMetric/statsd-http-proxy)
[![Code Climate](https://codeclimate.com/github/GoMetric/statsd-http-proxy/badges/gpa.svg?1)](https://codeclimate.com/github/GoMetric/statsd-http-proxy)
[![Coverage Status](https://coveralls.io/repos/github/GoMetric/statsd-http-proxy/badge.svg?branch=master)](https://coveralls.io/github/GoMetric/statsd-http-proxy?branch=master)

StatsD uses UDP connections, and  can not be used directly from browser. This server is a HTTP proxy to StatsD, useful for sending metrics to StatsD from frontend by AJAX.

Requests may be optionally authenticated using JWT tokens.

## Table of contents

* [Installation](#installation)
* [Requirements](#requirements)
* [Proxy client for browser](#proxy-client-for-browser)
* [Nginx config](#nginx-config)
* [Usage](#usage)
* [Authentication](#authentication)
* [Rest resources](#rest-resources)
  * [Heartbeat](#heartbeat)
  * [Count](#count)
  * [Gauge](#gauge)
  * [Timing](#timing)
  * [Set](#set)
* [Response](#response)
* [Testing](#testing)
* [Benchmark](#benchmark)
* [Useful resources](#useful-resources)


## Installation

```
git clone git@github.com:GoMetric/statsd-http-proxy.git
make build
```

Also available [Docker image](https://hub.docker.com/r/gometric/statsd-http-proxy/):

[![docker](https://img.shields.io/docker/pulls/gometric/statsd-http-proxy.svg?style=flat)](https://hub.docker.com/r/gometric/statsd-http-proxy/)

```
docker run -p 80:80 gometric/statsd-http-proxy:latest --verbose
```

Secure connection:

```
docker run -p 4433:4433 -v "$(pwd)":/certs/  gometric/statsd-http-proxy:latest --verbose --http-port=4433 --tls-cert=/certs/cert.pem --tls-key=/certs/key.pem
```

## Requirements

* [GoMetric/go-statsd-client](https://github.com/GoMetric/go-statsd-client) - StatsD client library for Go
* [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) - JSON Web Tokens builder and parser
* [gorilla/mux](https://github.com/gorilla/mux) - URL router and dispatcher

## Proxy client for browser

Basic implementation of proxy client may be found at https://github.com/GoMetric/statsd-http-proxy-client.

## Nginx config

Configuration of Nginx balancer:

```
server {
    listen 443 http2;

    server_name statsd-proxy.example.com;

    ssl on;
    ssl_certificate     /etc/pki/nginx/ssl.crt;
    ssl_certificate_key /etc/pki/nginx/ssl.key;

    upstream statsd_proxy {
        keepalive 100;
        server statsd-proxy-1:8825 max_fails=0;
        server statsd-proxy-2:8825 max_fails=0;
    }
    
    location / {
        proxy_pass http://statsd_proxy;
        proxy_redirect off;
        proxy_http_version 1.1;
        proxy_set_header Connection "keep-alive";
    }
}
```


## Usage

* Run server (HTTP):

```bash
statsd-http-proxy \
    --verbose \
    --http-host=127.0.0.1 \
    --http-port=8080 \
    --statsd-host=127.0.0.1 \
    --statsd-port=8125 \
    --jwt-secret=somesecret \
    --metric-prefix=prefix.subprefix
```

* Run server (HTTPS):

```bash
statsd-http-proxy \
    --verbose \
    --http-host=127.0.0.1 \
    --http-port=433 \
    --tls-cert=cert.pem \
    --tls-key=key.pem \
    --statsd-host=127.0.0.1 \
    --statsd-port=8125 \
    --jwt-secret=somesecret \
    --metric-prefix=prefix.subprefix
```

Print server version and exit:

```bash
statsd-http-proxy --version
```

Command line arguments:

| Parameter       | Description                          | Default value                                                                     |
|-----------------|--------------------------------------|-----------------------------------------------------------------------------------|
| verbose         | Print debug info to stderr           | Optional. Default false                                                           |
| http-host       | Host of HTTP server                  | Optional. Default 127.0.0.1. To accept connections on any interface, set to ""    |
| http-port       | Port of HTTP server                  | Optional. Default 80                                                              |
| http-timeout-read | The maximum duration in seconds for reading the entire request, including the body | Optional. Defaults to 1 second |
| http-timeout-write | The maximum duration in seconds before timing out writes of the respons | Optional. Defaults to 1 second  |
| http-timeout-idle | The maximum amount of time in seconds to wait for the next request when keep-alives are enabled | Optional. Defaults to 1 second |
| tls-cert        | TLS certificate for the HTTPS        | Optional. Default "" to use HTTP. If both tls-cert and tls-key set, HTTPS is used |
| tls-key         | TLS private key for the HTTPS        | Optional. Default "" to use HTTP. If both tls-cert and tls-key set, HTTPS is used |
| statsd-host     | Host of StatsD instance              | Optional. Default 127.0.0.1                                                       |
| statsd-port     | Port of StatsD instance              | Optional. Default 8125                                                            |
| jwt-secret      | JWT token secret                     | Optional. If not set, server accepts all connections                              |
| metric-prefix   | Prefix, added to any metric name     | Optional. If not set, do not add prefix                                           |
| version         | Print version of server and exit     | Optional                                                                          |

Sample code to send metric in browser with JWT token in header:

```javascript
$.ajax({
    url: 'http://127.0.0.1:8080/count/some.key.name',
    method: 'POST',
    headers: {
        'X-JWT-Token': 'some-jwt-token'
    },
    data: {
        value: 100500
    }
});
```

## Authentication

Authentication is optional. It based on passing JWT token to server, encrypted with secret, specified in `jwt-secret`
command line argument. If secret not configured in `jwt-secret`, then requests to server accepted without authentication.
Token sends to server in `X-JWT-Token` header or in `token` query parameter.

We recommend to use JWT tokens to prevent flood requests, setup JWT token expiration time, and update active JWT token in browser.

## Rest resources

See [statsd ductmentation](https://github.com/statsd/statsd/blob/master/docs/metric_types.md) about supported types.

### Heartbeat
```
GET /heartbeat
```
If server working, it responds with `OK`

### Count
```
POST /count/{key}
X-JWT-Token: {tokenString}
value=1&sampleRate=1
```

| Parameter  | Description                          | Default value                      |
|------------|--------------------------------------|------------------------------------|
| value      | Value. Negative to decrease          | Optional. Default 1                |
| sampleRate | Sample rate to skip metrics          | Optional. Default to 1: accept all |

### Gauge

Gauge is an arbitrary value. 
Only the last value during a flush interval is flushed to the backend. If the gauge is not updated at the next flush, it will send the previous value.
Gauge also may be set relatively to previously stored value. 
Is `shift` not passed, then `value` used. 
If `value` not passed, used default value equals 1.
To set a gauge to a negative number you need first set it to 0.

Absolute value:
```
POST /gauge/{key}
X-JWT-Token: {tokenString}
value=1
```

Shift of previous value:
```
POST /gauge/{key}
X-JWT-Token: {tokenString}
shift=-1
```

| Parameter  | Description                                     | Default value                      |
|------------|-------------------------------------------------|------------------------------------|
| value      | Integer value                                   | Optional. Default 1                |
| shift      | Signed int, relative to previously stored value | Optional                           |

### Timing
```
POST /timing/{key}
X-JWT-Token: {tokenString}
time=1234567&sampleRate=1
```

| Parameter  | Description                                   | Default value                      |
|------------|-----------------------------------------------|------------------------------------|
| time       | Time in milliseconds                          | Required                           |
| sampleRate | Float sample rate to skip metrics from 0 to 1 | Optional. Default to 1: accept all |

### Set
```
POST /set/{key}
X-JWT-Token: {tokenString}
value=1
```

| Parameter  | Description                          | Default value                      |
|------------|--------------------------------------|------------------------------------|
| value      | Integer value                        | Optional. Default 1                |

## Response

Server sends `200 OK` if send success, even StatsD server is down.

Other HTTP status codes:

| CODE             | Description                             |
|------------------|-----------------------------------------|
| 400 Bad Request  | Invalid parameters specified            |
| 401 Unauthorized | Token not sent                          |
| 403 Forbidden    | Token invalid/expired                   |
| 404 Not found    | Invalid url requested                   |
| 405 Wrong method | Request method not allowed for resource |

## Testing

It is useful for testing to start `netcat` UDP server,
listening for connections and watch incoming metrics.
To start server run:
```
nc -kluv localhost 8125
```
To send comamnd to statsd, run:

```
echo "counters" | nc localhost 8125
```

## Benchmark

Machine for benchmarking:

```
Intel(R) Core(TM) i5-2450M CPU @ 2.50GHz Dual Core / 8 GB RAM
```

Os:

```
Linux hp 4.15.0-65-generic #74-Ubuntu SMP Tue Sep 17 17:06:04 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux
```

Sysctl:

```
sudo sysctl net.ipv4.ip_local_port_range="15000 61000"
sudo sysctl net.ipv4.tcp_fin_timeout=30
```

Proxy's command line parameters:

```
$ GOMAXPROCS=2 ./bin/statsd-http-proxy --verbose --http-host=127.0.0.1 --http-port=8080 --statsd-host=127.0.0.1 --statsd-port=8125 --jwt-secret=somesecret
```

Siege test of Gorilla MUX Router without keep alive.

```
time siege -R <(echo connection = close) -c 255 -r 255 -H 'Connection: keep-alive' -H 'X-JWT-Token:eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzdGF0c2QtcmVzdC1zZXJ2ZXIiLCJpYXQiOjE1MDY5NzI1ODAsImV4cCI6MTg4NTY2Mzc4MCwiYXVkIjoiaHR0cHM6Ly9naXRodWIuY29tL3Nva2lsL3N0YXRzZC1yZXN0LXNlcnZlciIsInN1YiI6InNva2lsIn0.sOb0ccRBnN1u9IP2jhJrcNod14G5t-jMHNb_fsWov5c' "http://127.0.0.1:8080/count/a.b.c.d POST value=42"
** SIEGE 4.0.4
** Preparing 255 concurrent users for battle.
The server is now under siege...
Transactions:                  65025 hits
Availability:                 100.00 %
Elapsed time:                  64.70 secs
Data transferred:               0.00 MB
Response time:                  0.25 secs
Transaction rate:            1005.02 trans/sec
Throughput:                     0.00 MB/sec
Concurrency:                  252.40
Successful transactions:       65025
Failed transactions:               0
Longest transaction:            2.75
Shortest transaction:           0.00

real    1m4,728s
user    0m5,363s
sys     3m37,614s
```
Siege test of Gorilla MUX Router with keep alive.

```
time siege -R <(echo connection = keep-alive) -c 255 -r 255 -H 'Connection: keep-alive' -H 'X-JWT-Token:eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzdGF0c2QtcmVzdC1zZXJ2ZXIiLCJpYXQiOjE1MDY5NzI1ODAsImV4cCI6MTg4NTY2Mzc4MCwiYXVkIjoiaHR0cHM6Ly9naXRodWIuY29tL3Nva2lsL3N0YXRzZC1yZXN0LXNlcnZlciIsInN1YiI6InNva2lsIn0.sOb0ccRBnN1u9IP2jhJrcNod14G5t-jMHNb_fsWov5c' "http://127.0.0.1:8080/count/a.b.c.d POST value=42"
** SIEGE 4.0.4
** Preparing 255 concurrent users for battle.
The server is now under siege...
Transactions:                  65025 hits
Availability:                 100.00 %
Elapsed time:                   6.65 secs
Data transferred:               0.00 MB
Response time:                  0.02 secs
Transaction rate:            9778.20 trans/sec
Throughput:                     0.00 MB/sec
Concurrency:                  233.50
Successful transactions:       65025
Failed transactions:               0
Longest transaction:            0.72
Shortest transaction:           0.00
 

real    0m6,683s
user    0m5,089s
sys     0m4,324s
```

## Useful resources
* [https://github.com/etsy/statsd](https://github.com/etsy/statsd) - StatsD sources
* [Official Docker image for Graphite](https://github.com/graphite-project/docker-graphite-statsd)
* [Docker image with StatsD, Graphite, Grafana 2 and a Kamon Dashboard](https://github.com/kamon-io/docker-grafana-graphite)
* [Online JWT generator](http://jwtbuilder.jamiekurtz.com/)
* [Client for StatsD (Golang)](https://github.com/GoMetric/go-statsd-client)
