FROM golang:1.8.3-alpine

# Get git
RUN apk add --no-cache curl git

# Get glide
RUN go get github.com/Masterminds/glide

# Where factom-cli sources will live
WORKDIR $GOPATH/src/github.com/FactomProject/factom-cli

# Populate the source
COPY . .

# Install dependencies
RUN glide install -v

ARG GOOS=linux

# Build and install factom-cli
RUN go install

ENTRYPOINT ["/go/bin/factom-cli"]
