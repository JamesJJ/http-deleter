ARG BASE_IMAGE=gcr.io/distroless/base

ARG BUILD_BASE_IMAGE=golang:1.12-stretch

#-----------------------------------------#

FROM $BASE_IMAGE:debug as os_config_image

ARG BASE_IMAGE

ARG BUILD_BASE_IMAGE

RUN ["adduser","-h","/","-s","/sbin/nologin","-D","-H","app_user"]

#-----------------------------------------#

FROM $BUILD_BASE_IMAGE as build_image

ARG BASE_IMAGE

ARG BUILD_BASE_IMAGE

WORKDIR /go/src/app

COPY *.go ./

RUN go get -d -v ./...

RUN go test -v -cover ./...

RUN go install -v ./...

#-----------------------------------------#

FROM $BASE_IMAGE as service_image

ARG BASE_IMAGE

LABEL base_image=$BASE_IMAGE

LABEL owner_team=OPS

COPY --from=os_config_image --chown=root:root /etc/passwd /etc/group /etc/

COPY --from=build_image --chown=root:root /go/bin/app /app

USER app_user:app_user

ENTRYPOINT ["/app"]

