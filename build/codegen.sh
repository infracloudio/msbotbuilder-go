#!/bin/bash

OUT_DIR="pkg"

echo 'Generating code with openapi-generator-cli'
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/protocol/botframework.json -o /local/${OUT_DIR} -g go -c /local/protocol/openapi-config.yaml

# Fix permissions
sudo find ${OUT_DIR} -type d -exec chmod 755 {} \;
sudo find ${OUT_DIR} -type f -exec chmod 664 {} \;
