# Database Connection Manager
Local database connection manager.

## Requirements
- Go 1.26

## Usage
1. Build binary for command (file will be generated under `./bin/dbcm` by default).
```bash
$ mise build
```

2. (optional) Add path to bin file.
- if not set, call the command directly by its path, eg. `./bin/dbcm`.
```bash
$ export PATH=$PWD/bin:$PATH
```

3. Use the command to manage local database connections.
- compatible with: postgresql, mysql
```bash
# for the first time, initialize file to store managing states
$ dbcm init

# add connection to management
$ dbcm add psql <name> -h localhost -d testdb -u user -w pw

# connect to database by given name
$ dbcm connect <name>

# show connections
$ dbcm ls
```

## Run on docker environment
1. Build the environment.
```bash
$ mise docker-test
```

2. Use the `dbcm-linux` (linux compatible binary) command.
