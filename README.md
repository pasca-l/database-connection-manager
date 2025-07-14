# Database Connection Manager
Local database connection manager.

## Requirements
- Go 1.24

## Usage
1. Build binary for command (file will be generated under `./bin/dbcm` by default).
```bash
$ make build
```

2. Use the command to manage local database connections.
- compatible with: postgresql, mysql
```bash
# add connection to management
$ ./bin/dbcm add test psql -h localhost -d testdb -u user -w pw

# connect to database by given name
$ ./bin/dbcm connect test

# show connections and sessions
$ ./bin/dbcm ls
$ ./bin/dbcm sessions
```
