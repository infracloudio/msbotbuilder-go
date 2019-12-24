#!/bin/bash

set -o errexit
set -o nounset

OUT_DIR="."
uid="$(id -u)"
gid="$(id -g)"

echo 'Generating code with openapi-generator-cli'
sudo chown ${uid}:${gid} -R .

docker run --rm -v ${PWD}:/local --user ${uid}:${gid} \
     openapitools/openapi-generator-cli generate \
    -i /local/protocol/botframework.json \
    -o /local/${OUT_DIR} \
    -g go-server \
    -c /local/protocol/openapi-config.yaml
    
# Fix permissions
#sudo find ${OUT_DIR} -type d -exec chmod 775 {} \;
#sudo find ${OUT_DIR} -type f -exec chmod 664 {} \;
