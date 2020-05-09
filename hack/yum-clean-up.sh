#!/bin/bash
#
# Make sure package manager cache is cleaned up.
#

set -eu

rm -rv /var/cache /var/log/dnf* /var/log/yum.* || true
