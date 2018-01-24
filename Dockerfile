FROM scratch
ADD main.exe /
ADD data/config.yaml /data/
ENV PORT=80
ENTRYPOINT ["/main.exe", "--config",  "./data/config.yaml"]
EXPOSE 80
