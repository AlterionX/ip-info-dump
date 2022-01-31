# ip-info-dump

A simple serverless function that pulls info together from various locations and spits it back up in a simple format.

# Project structure

There are two directories of note: `env` and `src`

`env` should contain all files relevant to deployment while `src` contains all files relevant to the actual code and its tests.

Shell scripts to facilitate testing and deployment are in the root directory.

# Setup

There are two files that contain environment variables: `prod.tfvars` and `local.env.sh`.

The first is a var file for terraform to understand any environment variables that need to be loaded.

The second is a shell script to load any environment variables that need to be loaded locally.

The only environment variable that is relevant is `IPDUMP_VT_KEY`, which is a VirusTotal API key.

There are example versions of both (`env/prod.example.tfvars` and `local.example.env.sh`) in the repo. For security, the actual
files have been added to the `.gitignore`.

Once created, and the local environment set up for AWS, one can simply run `./deploy.sh` and `./test.sh` to either (re)deploy the
code to AWS or to test the code locally, respectively. `test.sh` will also provide a coverage report via `go test`'s `-cover`
option.

You do not need to be in the root directory to run these scripts.
