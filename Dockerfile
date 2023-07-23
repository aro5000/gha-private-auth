FROM golang:1.20-alpine as build
WORKDIR /app
ADD . .
ENV GOPATH /go
ENV CGO_ENABLED=0

RUN GOOS=linux GOARCH=amd64 go build

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gha-private-auth /gha-private-auth
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]