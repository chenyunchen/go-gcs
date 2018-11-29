# Building stage
FROM golang:1.11-alpine3.7

WORKDIR /go-gcs

RUN apk add --no-cache protobuf ca-certificates make git

# Source code, building tools and dependences
COPY src /go-gcs
COPY Makefile /go/src/github.com/linkernetworks/vortex

ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Taipei"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata

RUN make src.build
RUN mv build/src/cmd/filemanager/filemanager /go/bin

# Production stage
FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /go-gcs

# copy the go binaries from the building stage
COPY --from=0 /go/bin /go/bin

EXPOSE 7890
ENTRYPOINT ["/go/bin/filemanager", "-port", "7890", "-config", "/go-gcs/config/develop.json"]
