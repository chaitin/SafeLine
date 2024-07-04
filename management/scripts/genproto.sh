#!/usr/bin/env bash
#
# Generate all protobuf bindings.
# Run from repository root.

set -u

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
	echo "must be run from repository root"
	exit 255
fi

DIRS="webserver/proto tcontrollerd/proto"

echo "generating code"
protoc --version
for dir in ${DIRS}; do
	pushd "${dir}" || return
	  find . -type d -print0 | while IFS= read -r -d '' sdir ; do
	    pushd "${sdir}" || return
	      # shellcheck disable=SC2010
	      FS=$(ls | grep "\.proto\$")
        if [ -n "${FS}" ] ; then
          protoc --go_out=. --go_opt=paths=source_relative \
              --go-grpc_out=. --go-grpc_opt=paths=source_relative \
              "${FS}"

          goimports -local chaitin.cn -w ./*.pb.go
        fi
	    popd || return
	  done
	popd || return
done
