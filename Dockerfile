#
# Builder
#

FROM registry.fedoraproject.org/fedora:latest AS builder

ENV BUILDAH_ISOLATION="chroot"

COPY . /src
WORKDIR /src

RUN /src/hack/golang.sh && \
    /src/hack/deps.sh && \
    /src/hack/build-deps.sh

RUN make

#
# Application
#

FROM registry.fedoraproject.org/fedora:latest

COPY hack /src/hack

RUN /src/hack/storage.sh
ADD etc/containers/containers.conf /etc/containers/containers.conf
ADD etc/containers/storage.conf /etc/containers/storage.conf

RUN /src/hack/deps.sh && \
    /src/hack/build-deps.sh && \
    /src/hack/yum-clean-up.sh && \
    rm -rf /src

COPY --from=builder /src/_output/imagenie /usr/loca/bin/imagenie

ENTRYPOINT [ "/usr/loca/bin/imagenie" ]
