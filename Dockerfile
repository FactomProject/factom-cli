
FROM golang:1.15

# Where factom-cli sources will live
WORKDIR $GOPATH/src/github.com/FactomProject/factom-cli

# Populate the rest of the source
COPY . .

ARG GOOS=linux

# Build and install factom-cli
RUN make install

ENTRYPOINT ["/go/bin/factom-cli"]
