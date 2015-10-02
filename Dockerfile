FROM golang:1.5.1
EXPOSE 80
WORKDIR /go/src/github.com/manifest-destiny/api
COPY . ./
ENV GO15VENDOREXPERIMENT 1
RUN go build -o manifest-destiny-api .
CMD ["./manifest-destiny-api"]
