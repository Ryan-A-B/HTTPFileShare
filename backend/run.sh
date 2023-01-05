#!/bin/bash
set -ex

docker run --rm -it \
    --network host \
    -p 9000:9000 \
    http-file-share_backend
