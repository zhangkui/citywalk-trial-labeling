#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="citywalk-benzhi"
TARGET_PLATFORM="linux/amd64"

docker build --platform "${TARGET_PLATFORM}" -f benzhi.Dockerfile -t "${IMAGE_NAME}" .
