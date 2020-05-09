#!/bin/bash
#
# Installs the runtime dependencies.
#

set -eu

yum install -y \
    containers-common \
    fuse-overlayfs
