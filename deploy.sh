#! /bin/bash

set -euxo pipefail;

cd src;
env GOOS=linux GOARCH=amd64 go build -o ../bin/main;
cd ../bin;
zip -j main.zip main;
cd ../env;
terraform apply -var-file prod.tfvars;
