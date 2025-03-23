FROM golang:1.19 AS builder
MAINTAINER jinyashuo <jinyashuo@teleinfo.cn>
WORKDIR /opt/cndidresolve
COPY .. /opt/cndidresolve
ENV GO111MODULE=on
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cndidresolve .
RUN go clean -modcache

FROM scratch
MAINTAINER jinyashuo <jinyashuo@teleinfo.cn>
COPY --from=builder /opt/cndidresolve/cndidresolve /cndidresolve
EXPOSE 8080
ENTRYPOINT ["/cndidresolve"]