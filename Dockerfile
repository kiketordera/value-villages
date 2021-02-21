FROM golang:alpine as builder

RUN apk update && apk add git && apk add ca-certificates

COPY . $GOPATH/src/github.com/kiketordera/value-villages/
WORKDIR $GOPATH/src/github.com/kiketordera/value-villages/

RUN go get -d -v $GOPATH/src/github.com/kiketordera/value-villages/cmd/valuevillages
# For RaspberryPI
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/valuevillages $GOPATH/src/github.com/kiketordera/value-villages/cmd/valuevillages
# For Cloud Server
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/valuevillages $GOPATH/src/github.com/kiketordera/value-villages/cmd/valuevillages

FROM scratch
COPY --from=builder /go/bin/valuevillages /valuevillages
COPY --from=builder /go/src/github.com/kiketordera/value-villages/web/ /go/src/github.com/kiketordera/value-villages/web/
COPY --from=builder /go/src/github.com/kiketordera/value-villages/data-for-new-start/ /go/src/github.com/kiketordera/value-villages/data-for-new-start/
# If we also want to copy the Data of the Database:
# COPY --from=builder /go/src/github.com/kiketordera/value-villages/local-resources/ /go/src/github.com/kiketordera/value-villages/local-resources/

EXPOSE 8080/tcp

ENV GOPATH /go
ENTRYPOINT ["/valuevillages"]

