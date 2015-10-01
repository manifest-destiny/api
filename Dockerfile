FROM golang:1.5.1
MAINTAINER baydodd@gmail.com
ENV DIR /app/src/github.com/manifest-destiny/api
RUN mkdir -p $DIR
ADD . $DIR
WORKDIR $DIR
ENV GOPATH $GOPATH:/app
ENV GO15VENDOREXPERIMENT 1
RUN go build -o manifest-dstiny-api .
# RUN eval "$(./env)"
CMD ["./manifest-dstiny-api"]