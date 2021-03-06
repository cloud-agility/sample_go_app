FROM golang:1.8.5

WORKDIR /go/src/sample_go_app
COPY . .
RUN go get -d -v ./...
RUN go get -t ./...
RUN go install -v ./...
RUN go build

EXPOSE 8080
CMD ["go-wrapper","run"]
