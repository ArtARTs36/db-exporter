# db-exporter

db-exporter - simple app for export db schema to formats:
* Markdown `md`
* Class diagram `diagram`
* Go structures with db tags `go-structs`
* Goose migrations `goose`
* Migrations for [sql-migrate](https://github.com/rubenv/sql-migrate) `go-sql-migrate`
* Raw SQL Laravel migrations `laravel-migrations-raw`
* YAML fixtures `yaml-fixtures`

Supported database: PostgreSQL

usage:
```text
db-exporter driver-name dsn format out-dir [--table-per-file] [--with-diagram] [--without-migrations-table] [--tables=<value>] [--package=<value>] [--file-prefix=<value>]
[--commit-message=<value>] [--commit-push] [--commit-author=<value>] [--stat] [--debug]

Arguments
  driver-name                database driver name, required, available values: [pg]
  dsn                        data source name, required
  format                     exporting format, required, available values: [md, diagram, go-structs, goose, go-sql-migrate, laravel-migrations-raw, goose-fixtures]
  out-dir                    Output directory, required

Options
  table-per-file             Export one table to one file
  with-diagram               Export with diagram (only markdown)
  without-migrations-table   Export without migrations table
  tables                     Table list for export, separator: ","
  package                    Package name for code gen, e.g: models
  file-prefix                Prefix for generated files
  commit-message             Add commit with generated files and your message
  commit-push                Push commit with generated files
  stat                       Print stat
  import                     Import data from exported files
```

**Export from postgres to markdown**

```db-exporter pg "host=postgres user=root password=root dbname=cars" md ./docs```

## Using custom templates

[Twig syntax](https://twig.symfony.com) is used to compile templates. The Twig port is a [Stick](https://github.com/tyler-sommer/stick).

| Exporter       | Template                     | Description                                                  |
|----------------|------------------------------|--------------------------------------------------------------|
| md             | md/single-tables.md          | Template for generate single markdown file                   |
| md             | md/per-index.md              | Template for generate index markdown file (--table-per-file) |
| md             | md/per-table.tmd             | Template for generate table markdown file (--table-per-file) |
| diagram        | diagram/table.html           | Template for generate table                                  |
| go-structs     | go-structs/model.go.tpl      | Template for generate table                                  |
| goose          | goose/migration.sql          | Template for generate migration                              |
| goose-fixtures | goose/migration.sql          | Template for generate migration with fixtures                |
| go-sql-migrate | go-sql-migrate/migration.sql | Template for generate migration                              |
| laravel        | laravel/migration-raw.php    | Template for generate migration                              |
| grpc-crud      | grpc-crud/gprc.proto         | Template for generate protobuf                               |

You can download templates from [/templates](./templates)

In order for the db-exporter to use **your** templates, you need to place them in the `./db-exporter-templates` folder

## Use with GitHub Actions

You can run `db-exporter` as a GitHub action as follows:

```yaml
name: Generate documentation

permissions: write-all

on:
  push:
    branches:
      - master

jobs:
  generate-docs:
    services:
      postgres:
        image: postgres:12
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: cars
        ports:
          - 5499:5432

    runs-on: ubuntu-latest
    steps:
      - name: Checkout Source
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4 # action page: <https://github.com/actions/setup-go>
        with:
          go-version: 1.21.0

      - name: Install migrator
        run: go install github.com/pressly/goose/v3/cmd/goose@latest

      - name: Run migrations
        env:
          GOOSE_MIGRATION_DIR: './migrations'
        run: goose postgres "host=localhost port=5499 user=test password=test dbname=cars sslmode=disable" up

      - name: Generate markdown docs
        uses: artarts36/db-exporter@master
        with:
          driver-name: pg
          dsn: "host=localhost port=5499 user=test password=test dbname=cars sslmode=disable"
          format: md
          out-dir: ./docs
          commit-message: "chore: generate documentation for database schema"
          commit-push: true
          without-migrations-table: true
          with-diagram: true
````
