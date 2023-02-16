## About

This script allows backing up multiple file paths (specified in a config file) at once to an S3-compatible storage using Restic.

## How to use

Prepare `config.yaml` and `secrets.yaml`. (Examples are present in the repo.)

Then run:
```bash
restic-wrapper init  # should be run only once to prepare a remote repository
restic-wrapper  # to create a new snapshot
```