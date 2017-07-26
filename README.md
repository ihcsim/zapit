# url-scanner

[ ![Codeship Status for ihcsim/url-scanner](https://app.codeship.com/projects/52115f30-53eb-0135-fd18-160627fc0fd3/status?branch=master)](https://app.codeship.com/projects/235123)

## Table of Content

* [Problem Description](#problem-description)
* [Scaling Strategy](#scaling-strategy)
* [Development](#development)

## Problem Description
We have an HTTP proxy that is scanning traffic looking for malware URLs. Before allowing HTTP connections to be made, this proxy asks a service that maintains several databases of malware URLs if the resource being requested is known to contain malware.

Write a small web service, in the language/framework your choice, that responds to GET requests where the caller passes in a URL and the service responds with some information about that URL. The GET requests look like this:

```
GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```

The caller wants to know if it is safe to access that URL or not. As the implementer you get to choose the response format and structure. These lookups are blocking users from accessing the URL until the caller receives a response from your service.

## Scaling Strategy
Give some thoughts to the following:

_**The size of the URL list could grow infinitely, how might you scale this beyond the memory capacity of this VM? Bonus if you implement this.**_

One approach is to consider distributing the data across multiple nodes via sharding. Performance penalty will be incurred as additional time will be needed to look up the shard that hosts the requested data. We can also consider other storage solutions such as NFS or AWS EBS where we map our containers' data volumes to particular host paths, which are mounted to these external storage solutions.

In addition, we can also define a data retention policy where data entries are purged from the database when they satisfy certain criteria. Examples of data expiry criteria may include after some period of time, an URL's domain no longer exists etc.

If disk space is a concern, we can also investigate into data compression.

_**The number of requests may exceed the capacity of this VM, how might you solve that? Bonus if you implement this.**_

We can consider adding a load balancer in front of our service, and rely on a scheduler to scale our service's containers. This approach depends on the scheduler to provide the proxy abstraction to route the traffic to the container replicas.

_**What are some strategies you might use to update the service with new URLs? Updates may be as much as 5 thousand URLs a day with updates arriving every 10 minutes.**_

One approach is to consider hosting a canonical data source where all database instances pull their data from. New updates are applied to this source. They can be either pushed to or pulled by all existing databases throughout the day, at different time intervals. A naive implementation is one where the canonical data source are some data files stored in an S3 bucket. At different times throughout the day, our service replicas will download this file, delete their respective database, and re-populate their database with entries in these files.

## Development
To get the source:
```
$ go get github.com/ihcsim/url-scanner
```

To run the tests:
```
$ go test -v -cover -race ./...
```

To build the server:
```
$ go build -v github.com/ihcsim/url-scanner/cmd/server/...
```

To build the Docker image:
```
$ docker image build --rm -t <image_tag> .
```
