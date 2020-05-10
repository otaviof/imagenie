#!/bin/bash
#
# Capturing "buildRequires" from rpm-spec:
#   https://github.com/containers/buildah/blob/master/contrib/rpm/buildah.spec
#

set -eu

yum install -y \
    btrfs-progs-devel \
    device-mapper-devel \
    git \
    glib2-devel \
    gpgme-devel \
    libassuan-devel \
    libseccomp-devel \
    make
