# factom-cli Docker Helper

The factom-cli Docker Helper is a simple tool to help build and run factom-cli as a container

## Prerequisites

You must have at least Docker v17 installed on your system.

Having this repo cloned helps too 😇

## Build
From wherever you have cloned this repo, run

`docker build -t factom-cli_container .`

(yes, you can replace **factom-cli_container** with whatever you want to call the container.  e.g. **factom-cli**, **foo**, etc.)

#### Cross-Compile
To cross-compile for a different target, you can pass in a `build-arg` as so

`docker build -t factom-cli_container --build-arg GOOS=darwin .`

## Run
`docker run --rm factom-cli_container <some command here>`

e.g.

`docker run --rm factom-cli_container get heights`

**Note** - In the above, replace **factom-cli_container** with whatever you called it when you built it - e.g. **factom-cli**, **foo**, etc.


## Copy
So yeah, you want to get your binary _out_ of the container. To do so, you basically mount your target into the container, and copy the binary over, like so


`docker run --rm --entrypoint='' -v <FULLY_QUALIFIED_PATH_TO_TARGET_DIRECTORY>:/destination factom-cli_container /bin/cp /go/bin/factom-cli /destination`

e.g.

`docker run --rm --entrypoint='' -v /tmp:/destination factom-cli_container /bin/cp /go/bin/factom-cli /destination`

which will copy the binary to `/tmp/factom-cli`

**Note** : You should replace ** factom-cli_container** with whatever you called it in the **build** section above  e.g. **factom-cli**, **foo**, etc.

#### Cross-Compile
If you cross-compiled to a different target, your binary will be in `/go/bin/<target>/factom-cli`.  e.g. If you built with `--build-arg GOOS=darwin`, then you can copy out the binary with

`docker run --rm --entrypoint='' -v <FULLY_QUALIFIED_PATH_TO_TARGET_DIRECTORY>:/destination factom-cli_container /bin/cp /go/bin/darwin_amd64/factom-cli /destination`

e.g.

`docker run --rm --entrypoint='' -v /tmp:/destination factom-cli_container /bin/cp /go/bin/darwin_amd64/factom-cli /destination` 

which will copy the darwin_amd64 version of the binary to `/tmp/factom-cli`

**Note** : You should replace ** factom-cli_container** with whatever you called it in the **build** section above  e.g. **factom-cli**, **foo**, etc.
