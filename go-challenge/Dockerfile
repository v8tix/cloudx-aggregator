FROM golang:1.21.4-bullseye
RUN ["/bin/bash", "-c", "mkdir /src"]
WORKDIR /src
ADD . .
RUN ["/bin/bash", "-c", "GOOS=linux GOARCH=amd64 go build -buildvcs=false -o=/src/wsclient /src/cmd/app"]
RUN ["/bin/bash", "-c", "chmod +x /src/wsclient"]
ENTRYPOINT ["/src/wsclient"]

