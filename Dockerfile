FROM golang:1.11.13 AS build-env
WORKDIR /go/src/github.com/turbosonic/api-gateway
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app.exe app.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=build-env /go/src/github.com/turbosonic/api-gateway/app.exe .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV PORT=80
ENTRYPOINT ["/app.exe", "--config",  "./data/config.yaml"]
EXPOSE 80