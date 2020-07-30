FROM golang:alpine as builder

RUN apk update && apk add git 
COPY . $GOPATH/src/lacazethomas/audiApiExporter/
WORKDIR $GOPATH/src/lacazethomas/audiApiExporter/
RUN go get -d -v


RUN go build -o /go/bin/audiApiExporter


FROM alpine
EXPOSE 9158
COPY --from=builder /go/bin/audiApiExporter /bin/audiApiExporter
ENTRYPOINT ["/bin/audiApiExporter"]