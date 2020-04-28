#!/usr/bin/env bash
NAME=shadowsocks2-yee

go generate
for os in darwin freebsd windows; do
  echo build ${NAME}_${os}_i386
  GOOS=${os} GOARCH=386 go build -o releases/${NAME}_${os}_i386
  echo build ${NAME}_${os}_amd64
  GOOS=${os} GOARCH=amd64 go build -o releases/${NAME}_${os}_amd64
done
for arch in 386 amd64 arm arm64; do
  echo build ${NAME}_linux_${arch}
  GOOS=linux GOARCH=${arch} go build -o releases/${NAME}_linux_${arch}
done
