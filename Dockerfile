FROM golang:1.10 AS build
WORKDIR /go/src
COPY pkg ./pkg
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o botbuilder .

FROM scratch AS runtime
COPY --from=build /go/src/botbuilder ./
EXPOSE 8080/tcp
ENTRYPOINT ["./botbuilder"]
