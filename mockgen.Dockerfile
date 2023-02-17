FROM golang:1.15

RUN cd /tmp
RUN GO111MODULE=on go get github.com/golang/mock/mockgen@v1.6.0
RUN cd -

WORKDIR /src
ENTRYPOINT ["go", "generate", "-v", "./..."]
