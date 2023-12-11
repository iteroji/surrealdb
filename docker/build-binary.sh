#/bin/bash

set -eo pipefail

ROOT_DIR=$(dirname $(dirname $(realpath $0)))

: ${DOCKER_BINARY:=docker}

cd $ROOT_DIR

if [[ $(uname -m) == "x86_64" ]]; then
    export CARGO_EXTRA_FEATURES="$CARGO_EXTRA_FEATURES,storage-fdb"
    export DOCKER_BUILD_TARGET="builder-fdb" # Target for the binary builder
else
    echo "###"
    echo "###> WARNING!"
    echo "###> This architecture ('$(uname -m)') doesn't support FoundationDB."
    echo "###"

    export CARGO_EXTRA_FEATURES="http-compression,storage-tikv"
    export DOCKER_BUILD_TARGET="builder" # Target for the binary builder
fi

echo
echo "###"
echo "###> Build SurrealDB binary"
echo "###"
echo
$DOCKER_BINARY build -t surrealdb-local/builder --target $DOCKER_BUILD_TARGET -f docker/Dockerfile .
$DOCKER_BINARY run --rm -v $(pwd):/surrealdb -w /surrealdb surrealdb-local/builder cargo build --features ${CARGO_EXTRA_FEATURES}

echo "###"
echo "###> Build SurrealDB docker images"
echo "###"
