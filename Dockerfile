FROM golang:1.8.3

WORKDIR /go/src/github.com/ihcsim/zapit
COPY . .
RUN go install -v github.com/ihcsim/zapit/cmd/server
ENTRYPOINT ["server"]
