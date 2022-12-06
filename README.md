# gigi

## Specification

- Go 1.17

## Setup

```bash
# Install make, docker and write .env
$ bash script/setup.sh

# create docker container
$ make docker/run
```

## Make commands

```bash
# start container
$ make docker/start

# stop container
$ make docker/stop

# remove container
$ make docker/rm

# view logs
$ make docker/logs
```
