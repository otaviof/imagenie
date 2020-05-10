#
# Builder
#

ARG BUILDER_IMAGE="registry.fedoraproject.org/fedora:latest"
FROM $BUILDER_IMAGE AS builder

COPY . /src
WORKDIR /src

RUN make

#
# Application
#

FROM registry.fedoraproject.org/fedora:latest

COPY hack /src/hack

RUN /src/hack/storage.sh
ADD etc/containers/containers.conf /etc/containers/containers.conf
ADD etc/containers/storage.conf /etc/containers/storage.conf

VOLUME /var/lib/containers/storage
VOLUME /var/lib/shared

RUN /src/hack/deps.sh && \
    /src/hack/yum-clean-up.sh && \
    rm -rf /src

COPY --from=builder /src/_output/imagenie /usr/local/bin/imagenie

ENTRYPOINT [ "/usr/local/bin/imagenie" ]
