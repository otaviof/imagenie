#!/bin/bash

set -eux

IMAGES_DIR="/var/lib/shared/overlay-images"
LAYERS_DIR="/var/lib/shared/overlay-layers"

mkdir -p ${IMAGES_DIR} ${LAYERS_DIR}

touch ${IMAGES_DIR}/images.lock
touch ${LAYERS_DIR}/layers.lock