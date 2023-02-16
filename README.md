## About

This script allows backing up multiple file paths (specified in a config file) at once to an S3-compatible storage using Restic.

## TODO
- [ ] Allow working with multiple profiles (where `profile = repo + paths + password`)
- [ ] Allow working with non-S3 storages
- [ ] Support Windows
- [ ] Add pruning (like [here](https://pypi.org/project/runrestic/))

## How to build

```bash
go build .
```

## How to use

Prepare `config.yaml` and `secrets.yaml`. (Examples are present in the repo.)

Then run:
```bash
./restic-wrapper init  # should be run only once to prepare a remote repository
./restic-wrapper backup  # to create a new snapshot
./restic-wrapper snapshots  # to list created snapshots
./restic-wrapper help  # to print help
```