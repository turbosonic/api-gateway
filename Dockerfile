FROM scratch
ADD main.exe /
ADD config.yaml /
CMD ["/main.exe", "--config",  "/config.yaml"]
EXPOSE 8080