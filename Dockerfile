#
# Builder
#

FROM registry.fedoraproject.org/fedora:latest AS builder

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

VOLUME /var/lib/containers/storage
VOLUME /var/lib/shared

RUN /src/hack/deps.sh && \
    /src/hack/yum-clean-up.sh && \
    rm -rf /src

COPY --from=builder /src/_output/imagenie /usr/local/bin/imagenie

ENTRYPOINT [ "/usr/local/bin/imagenie" ]
