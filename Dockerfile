FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN go mod tidy
RUN go get github.com/markbates/pkger/cmd/pkger
RUN /go/bin/pkger
RUN go generate ./...
RUN go build -o ./bin/git-file-webserver .

FROM alpine:latest
RUN apk add git
COPY --from=builder /build/bin/git-file-webserver /bin/
VOLUME /config
WORKDIR /config
EXPOSE 80
CMD ["/bin/git-file-webserver"]
