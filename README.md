# Database Connection Manager
Local database connection manager.
> currently compatible with: PostgreSQL, MySQL

## Installation
### Homebrew
```bash
$ brew install --cask pasca-l/tap/database-connection-manager
```

### Local build (requires Go 1.26)
1. Build binary for command (file will be generated under `./bin/dbcm` by default).
```bash
$ mise build
```

2. (optional) Add path to bin file.
- if not set, call the command directly by its path, eg. `./bin/dbcm`.
```bash
$ export PATH=$PWD/bin:$PATH
```

## Usage
```bash
# for the first time, initialize file to store managing states
$ dbcm init

# add connection to management
# with psql syntax for PostgreSQL
$ dbcm add psql <name> -h localhost -d db -U user -w pw -p 5432
# with mysql syntax for MySQL
$ dbcm add mysql <name> -h localhost -D db -u user -p pw -P 3306

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
