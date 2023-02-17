## About

This script allows backing up multiple file paths (specified in a config file) at once to an S3-compatible storage using Restic.

## TODO
- [ ] Allow ignoring certain file paths
- [ ] Add help message
- [x] Allow working with multiple profiles (where `profile = repo + paths + password`)
    - *Currently, there's no way to specify a profile when running commands: the commands apply to all profiles at once.*
- [ ] Allow working with non-S3 storages
- [ ] Support Windows
- [ ] Add pruning (like [here](https://pypi.org/project/runrestic/))
- [ ] Clean up the output

## How to build

```bash
go build .
```

## How to use

- Prepare `config.yaml` and `secrets.yaml`. (Examples are present in the repo.)
- `config.yaml` should be either placed inside `/etc/restic-wrapper` or linked to by the `RESTIC_WRAPPER_CONFIG_PATH` environment variable.

Then run:
```bash
./restic-wrapper init  # should be run only once to prepare a remote repository
./restic-wrapper backup  # to create a new snapshot
./restic-wrapper snapshots  # to list created snapshots
./restic-wrapper help  # to print help
```