set -eu

useradd build

yum -y update
yum -y reinstall shadow-utils
yum -y install buildah fuse-overlayfs --exclude container-selinux

rm -rf /var/cache /var/log/dnf* /var/log/yum.*
