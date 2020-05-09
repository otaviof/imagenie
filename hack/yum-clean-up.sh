#!/bin/bash
#
# Make sure package manager cache is cleaned up.
#

set -eu

rm -rfv /var/cache /var/log/dnf* /var/log/yum.*
