###
# STAGE: fdb-static-libs
# This stage is used to statically build the foundationdb libraries
###
FROM --platform=linux/amd64 docker.io/foundationdb/build:centos7-latest AS fdb-static-libs

# Prepare foundationDB code
RUN git clone --depth 1 https://github.com/apple/foundationdb.git /foundationdb
COPY docker/files/fdb.patch /foundationdb/fdb.patch
RUN cd /foundationdb && git fetch --tags && git checkout 7.1.42 && git apply fdb.patch
RUN mkdir /build_output

# Build foundationdb static libraries
RUN source /opt/rh/devtoolset-11/enable && cmake -S /foundationdb -B /build_output -G Ninja
RUN source /opt/rh/devtoolset-11/enable && ninja -C /build_output fdb_c

###
# STAGE: builder
# This stage is used to build the SurrealDB binary
#
# It builds the SurrealDB Linux binary and makes it depend on GLIBC 2.18 thanks to using RockyLinux 8 as a base image.
###

FROM docker.io/rockylinux:8 as builder

RUN yum install -y cargo rust-1.71.* patch clang-16.0.* gcc-toolset-13 git

ENV CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER=aarch64-redhat-linux-gcc

###
# STAGE: builder-fdb
# This stage is used to build the SurrealDB binary with FoundationDB support
#
# It builds the SurrealDB Linux binary and makes it depend on GLIBC 2.18 thanks to using RockyLinux 8 as a base image.
###

FROM builder as builder-fdb

RUN sed -i 's/enabled=0/enabled=1/' /etc/yum.repos.d/Rocky-Devel.repo
RUN yum install -y openssl-static-1.1.1k

ENV FDB_CLIENT_STATIC=1
ENV FDB_CLIENT_LIB_PATHS=/usr/local/lib:/usr/lib64

COPY --from=fdb-static-libs /build_output/lib/* /usr/local/lib/

###
# Final Images
###

#
# Production image (CI environment)
#
FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/glibc-dynamic:latest as prod-ci

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG ARTIFACT_PREFIX

COPY ${ARTIFACT_PREFIX}.${TARGETOS}-${TARGETARCH}/surreal /surreal

ENTRYPOINT ["/surreal"]

#
# Development image (CI environment)
#
FROM --platform=$TARGETPLATFORM cgr.dev/chainguard/glibc-dynamic:latest-dev as dev-ci

ARG TARGETPLATFORM

USER root

COPY --from=prod-ci /surreal /surreal

ENTRYPOINT ["/surreal"]

#
# Production image
#
FROM cgr.dev/chainguard/glibc-dynamic:latest as prod

ARG SURREALDB_BINARY=target/release/surreal

COPY ${SURREALDB_BINARY} /surreal

ENTRYPOINT ["/surreal"]

#
# Development image
#
FROM cgr.dev/chainguard/glibc-dynamic:latest-dev as dev

USER root

COPY --from=prod /surreal /surreal

ENTRYPOINT ["/surreal"]
