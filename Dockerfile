FROM golang:1.8.3

WORKDIR /go/src/github.com/ihcsim/url-scanner
COPY . .
RUN go install -v github.com/ihcsim/url-scanner/cmd/server
ENTRYPOINT ["server"]
