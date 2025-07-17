# Database Connection Manager
Local database connection manager.

## Requirements
- Go 1.24

## Usage
1. Build binary for command (file will be generated under `./bin/dbcm` by default).
```bash
$ make build
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
$ dbcm add test psql -h localhost -d testdb -u user -w pw

# connect to database by given name
$ dbcm connect test

# show connections and sessions
$ dbcm ls
$ dbcm sessions
```

## Run on development mode
1. Set up docker containers.
```bash
$ make start
```

2. Enter docker container with database client.
```bash
$ make client
```

3. Use the `dbcm` command.
