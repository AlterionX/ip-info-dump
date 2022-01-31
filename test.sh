#! /bin/bash

set -euxo pipefail

# Load environment variables
source ./env/localenv.sh;
cd src;
go test ./... -cover;
