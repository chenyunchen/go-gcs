# Building stage
FROM golang:1.11-alpine3.7

WORKDIR /go-gcs

RUN apk add --no-cache protobuf ca-certificates make git

# Source code, building tools and dependences
COPY src /go-gcs/src
COPY Makefile /go-gcs
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod /go-gcs
COPY go.sum /go-gcs


ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Taipei"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata
# Force the go compiler to use modules
ENV GO111MODULE=on

RUN go mod download
RUN make src.build
RUN mv build/src/cmd/filemanager/file_manager /go/bin

# Production stage
FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /go-gcs

# copy the go binaries from the building stage
COPY --from=0 /go/bin /go/bin

# select the config file for deployment
ARG CONFIG=config/local.json
ENV CONFIG=/go-gcs/$CONFIG

EXPOSE 7890
ENTRYPOINT /go/bin/file_manager -port 7890 -config $CONFIG
