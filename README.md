# zapit

[ ![Codeship Status for ihcsim/url-scanner](https://app.codeship.com/projects/52115f30-53eb-0135-fd18-160627fc0fd3/status?branch=master)](https://app.codeship.com/projects/235123)

[![buddy pipeline](https://app.buddy.works/ihcsim/zapit/pipelines/pipeline/58248/badge.svg?token=a4e6f4e142f7c03b31711b30f8217b8915a52381c4ab1102a43ae0ab6d2ff978 "buddy pipeline")](https://app.buddy.works/ihcsim/zapit/pipelines/pipeline/58248)

zapit provides a scanner that checks a URL or IP to determine if the endpoint is on the ZeuS Tracker's blocklists.

## Table of Content

* [Introduction](#introduction)
* [Prerequisites](#prerequisites)
* [Request Format](#request-format)
* [Getting Started](#getting-started)
* [Development](#development)

## Introduction

zapit is made up of 4 components:

* Nginx proxies traffic into the system.
* Scanner performs the scan on the submitted endpoint.
* Redis stores a list of blocked endpoints obtained from the [ZeuS Tracker](https://zeustracker.abuse.ch/blocklist.php).
* Feeder polls the ZeuS Tracker website and RSS feed for new blocked URLs, at a configurable regular interval.

![System Design](https://github.com/ihcsim/zapit/raw/master/img/system-design.png)
An endpoint is marked as safe if it isn't found in zapit's database.

To counter the different permutations of paths, query strings, anchors and subdomains that can be added to masquerade a malicious server's hostname, zapit performs a two-pass scan on every endpoint it receives.

During the first pass, zapit strips away the endpoint's additional paths, query strings, anchors and subdomains in order to perform a scan on either the URL's second-level domain name or the IPv4 address. For example, given the URLs blog.example.com and support.eu.example.com, the example.com domain name will be scanned. The general idea is that if a domain is marked as unsafe, all its subdomains and paths will be unsafe.

If the first pass returns a positive result, indicating that the domain is safe, then a second pass is triggered. During the second pass, zapit scans the endpoint with its submitted subdomains and paths intact. If the endpoint passed the second scan, then it's marked as safe.

The feeder is scheduled to read from the ZeuS Tracker site every 30 minutes, and update the Redis database accordingly.

zapit reads the blocked lists from the following sites:

* https://zeustracker.abuse.ch/blocklist.php?download=baddomains - ZeuS domain blocklist "BadDomains"
* https://zeustracker.abuse.ch/blocklist.php?download=badips - ZeuS IP blocklist "BadIPs"
* https://zeustracker.abuse.ch/rss.php - This feed shows the latest twenty ZeuS hosts which the tracker has captured.

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
