#!/bin/bash
set -ex

docker run --rm -it \
    -p 3000:3000 \
    -v $(pwd):$(pwd) \
    -w $(pwd) \
    node:16 yarn start