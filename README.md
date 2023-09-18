# gigi

## Specification

- Go 1.20

## Setup

> **Note**  
> BOT requires `SERVER MEMBERS INTENT` and `MESSAGE CONTENT INTENT`.

```bash
# create docker container (build included)
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
