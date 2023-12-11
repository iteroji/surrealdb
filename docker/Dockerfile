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

RUN yum install -y cargo rust-1.71.* patch clang-16.0.*

###
# STAGE: builder-fdb
# This stage is used to build the SurrealDB binary with FoundationDB support
#
# It builds the SurrealDB Linux binary and makes it depend on GLIBC 2.18 thanks to using RockyLinux 8 as a base image.
###

FROM builder as builder-fdb

RUN sed -i 's/enabled=0/enabled=1/' /etc/yum.repos.d/Rocky-Devel.repo
RUN yum install -y openssl-static

ENV FDB_CLIENT_STATIC=1
ENV FDB_CLIENT_LIB_PATHS=/usr/local/lib:/usr/lib64

COPY --from=fdb-static-libs /build_output/lib/* /usr/local/lib/

###
# Final Images
###

#
# Development image
#
FROM cgr.dev/chainguard/glibc-dynamic:latest-dev as dev
ARG SURREALDB_BINARY=target/release/surreal

USER root
COPY ${SURREALDB_BINARY} /surreal

ENTRYPOINT ["/surreal"]

#
# Production image
#
FROM cgr.dev/chainguard/glibc-dynamic:latest as prod

ARG SURREALDB_BINARY=target/release/surreal

COPY ${SURREALDB_BINARY} /surreal

ENTRYPOINT ["/surreal"]
