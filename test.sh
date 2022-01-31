#! /bin/bash

set -euxo pipefail
cd "${0%/*}";

# Load environment variables
source ./env/local.env.sh;
cd src;
go test ./... -cover;
