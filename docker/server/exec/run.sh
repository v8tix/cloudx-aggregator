#!/usr/bin/env bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "${dir}/../configuration/const.sh"

docker run -d --rm --name "${IMAGE}" --hostname "${IMAGE}" -p "${SERVER_HOST_PORT}":"${SERVER_CONTAINER_PORT}" "${IMAGE_TAG}"