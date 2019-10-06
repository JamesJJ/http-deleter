#!/bin/bash

export BUILD_BASE_IMAGE='temp-image-ht-del-106'


docker run -it --rm -v "$(pwd)/:/mnt/" "golang:1.12-stretch" find /mnt -type f -not -path '*/.git/*' -name '*.go' -exec gofmt -w {} \;
if [ $? -ne 0 ]; then
  exit 1
fi


docker images -q "${BUILD_BASE_IMAGE}:cached" | grep -c '[[:alnum:]]' >/dev/null

if [ $? -ne 0 ]; then
  docker build --target=build_image -t "${BUILD_BASE_IMAGE}:cached" .
  if [ $? -ne 0 ]; then
    exit 1
  fi
fi

docker build -t "${BUILD_BASE_IMAGE}:app_build" --build-arg BUILD_BASE_IMAGE="${BUILD_BASE_IMAGE}:cached" . \
&& echo '= = = = =' \
&& docker run -it --rm -e "HTTP_DELETER_URLS=${HTTP_DELETER_URLS}" "${BUILD_BASE_IMAGE}:app_build"



