FROM golang:1.5.1
EXPOSE 80
ENV DIR github.com/manifest-destiny/api
WORKDIR $GOPATH/src/$DIR
COPY . ./
ENV GO15VENDOREXPERIMENT 1
RUN go build -o migrate $DIR/vendor/github.com/mattes/migrate
RUN chmod a+x ./migrate
RUN cd service && go build -o api .
CMD ["./service/api"]
