FROM golang:alpine as builder

RUN apk update && apk add git 
COPY . $GOPATH/src/lacazethomas/renaultApiExporter/
WORKDIR $GOPATH/src/lacazethomas/renaultApiExporter/
RUN go get -d -v


RUN go build -o /go/bin/renaultApiExporter


FROM alpine
EXPOSE 9158
COPY --from=builder /go/bin/renaultApiExporter /bin/renaultApiExporter
ENTRYPOINT ["/bin/renaultApiExporter"]