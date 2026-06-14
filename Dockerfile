FROM golang:1.26.4-trixie

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED="1"
