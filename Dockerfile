FROM scratch
ADD main.exe /
ADD web.v1.configuration.yaml /config.yaml
CMD ["/main.exe", "--config",  "/config.yaml"]