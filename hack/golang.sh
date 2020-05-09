#!/bin/bash
#
# Deploys latest Golang from "rawhide" repository
#

set -eu

yum install -y fedora-repos-rawhide
yum install -y --nogpgcheck --allowerasing --enablerepo=rawhide golang
