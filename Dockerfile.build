FROM golang:1.8 AS build
COPY . /go/src/docker-registry-cleaner

WORKDIR /go/src/docker-registry-cleaner
RUN go get -d ./... && \
 go get -t && \
 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o docker-registry-cleaner .

# copy the binary from the build stage to the final stage
FROM alpine:3.5
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build /go/src/docker-registry-cleaner/docker-registry-cleaner /docker-registry-cleaner
CMD ["/docker-registry-cleaner"]