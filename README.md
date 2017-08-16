# zapit

[ ![Codeship Status for ihcsim/url-scanner](https://app.codeship.com/projects/52115f30-53eb-0135-fd18-160627fc0fd3/status?branch=master)](https://app.codeship.com/projects/235123)

zapit provides a scanner that checks a URL or IP to determine if the URL is on the ZeuS Tracker's blocklists. It is made up of 4 components:

* nginx proxies traffic in and out of the system.
* scanner handles user's request by extracting the URL to be scanned from the request.
* redis stores a list of blocked URLs and IPs obtained from the [ZeuS Tracker](https://zeustracker.abuse.ch/blocklist.php).
* feeder polls the ZeuS Tracker website and RSS Feed for new blocked URLs, at a configurable regular interval.

![System Design](https://github.com/ihcsim/zapit/raw/master/img/system-design.png)

zapit reads from the following lists:

* https://zeustracker.abuse.ch/blocklist.php?download=baddomains - ZeuS domain blocklist "BadDomains"
* https://zeustracker.abuse.ch/blocklist.php?download=badips - ZeuS IP blocklist "BadIPs"
* https://zeustracker.abuse.ch/rss.php - This feed shows the latest twenty ZeuS hosts which the tracker has captured.
* https://zeustracker.abuse.ch/removals.php?show=rss - This feed shows the ZeuS hosts which were removedf from the tracker list.

## Table of Content

* [Prerequisites](#prerequisites)
* [System Design](#system-design)
* [Request Format](#request-format)
* [Getting Started](#getting-started)
* [Scaling Strategy](#scaling-strategy)
* [Development](#development)

## Prerequisites
The following is a list of software needed to run zapit:

* Docker 17.05 CE
* Docker Compose 1.13.0

## Request Format
The URL to be scanned will be appended as query string in a `GET` request as follows:
```
GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```

## Getting Started
Use Docker Compose to start the service:
```
$ docker-compose -p zapit -d up
```

Use `curl` to test the service:
```
$ curl localhost:8080/urlinfo/1/localhost
{"url":"localhost","isSafe":true}

$ curl localhost:8080/urlinfo/1/127.0.0.1
{"url":"127.0.0.1","isSafe":true}

$ curl localhost:8080/urlinfo/1/google.com
{"url":"google.com","isSafe":true}

$ curl localhost:8080/urlinfo/1/jjwire.biz
{"url":"jjwire.biz","isSafe":false}

$ curl localhost:8080/urlinfo/1/162.246.57.205
{"url":"162.246.57.205","isSafe":false}

$ curl localhost:8080/urlinfo/1/gmailsecurityteam.com?foo=bar&foo2=bar3
{"url":"gmailsecurityteam.com","isSafe":false}
```
Query strings must be URL-encoded.

The following configurations can be overridden with environment variables:

Variables      | Descriptions                            | Defaults
-------------- | --------------------------------------- | -------
`LB_PORT`      | TCP port that nginx listens on          | 8080
`SCANNER_PORT` | TCP port that the `scanner` listens on  | 8080
`REDIS_PORT`   | TCP port that the Redis listens on      | 6379
`DB_UPDATE_INTERVAL` | The time interval (in minutes) between polling the ZeuS web sites and RSS feed. Must satisfies the Go `time.Duration` format | 30m

The `.env` file contains defaults that docker-compose uses.

## Development
To get the source:
```
$ go get github.com/ihcsim/zapit
```

To run the tests:
```
$ go test -v -cover -race ./...
```

To build the server:
```
$ go build -v github.com/ihcsim/zapit/cmd/server/...
```

To build the Docker image:
```
$ docker image build --rm -t <image_tag> .
```
