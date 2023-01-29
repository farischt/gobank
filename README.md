# Gobank

This project is a simple bank CRUD API, written using standard http and sql packages.

## Requirements:

- Fill the config files following the given [examples](./config/config.development.example.json).

- **Make sure to update** your [docker-compose.yml](./docker-compose.yml) file to reflect your configurations.

- [Compile Daemon](https://github.com/githubnemo/CompileDaemon) (globaly installed) to run the server in development mode:

```bash
    go install github.com/githubnemo/CompileDaemon@latest
```

- [Golang-migrate](https://github.com/golang-migrate/migrate) (globaly installed) to run migrations.

## Start

<br>

Install all the dependencies:

```bash
    make deps
```

Start your postgres database:

```bash
    docker compose up -d
```

Start the web server in dev mode:

```bash
    make dev
```

## Migrations

<br>

- Create a migration:

```bash
    make migrate-create name=<YOUR_MIGRATION_NAME>
```

- Migrate up:

```bash
    make migrate-up
```

- Migrate down:

```bash
    make migrate-down
```

- Revert changes to a specific version:

```bash
    make migrate-goto v=<TARGET_VERSION>
```

- Fix a dirty version:

```bash
    make migrate-fix v=<TARGET_VERSION>
```
