#! /bin/bash

set -euxo pipefail

# Load environment variables
source ./env/local.env.sh;
cd src;
go test ./... -cover;
