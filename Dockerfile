FROM golang:1.18-buster AS builder

LABEL maintainer="Harri Avellan <harri@xavellan.tech>"

ARG VERSION=dev

WORKDIR /go/src/app
COPY main.go plugin.go go.mod go.sum /go/src/app/
RUN go build -o main -ldflags=-X=main.version=${VERSION} main.go plugin.go

FROM debian:buster-slim
RUN apt update && apt install -y ca-certificates
COPY --from=builder /go/src/app/main /go/bin/main
ENV PATH="/go/bin:${PATH}"
CMD ["main"]
