FROM golang:latest
WORKDIR /usr/local/bin
COPY ./network_scanner .
RUN chmod +x /usr/local/bin/network_scanner
ENTRYPOINT ["/usr/local/bin/network_scanner"]
