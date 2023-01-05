#!/bin/bash
set -ex

docker run --rm -it \
    -p 80:80 \
    --network host \
    http-file-share_frontend