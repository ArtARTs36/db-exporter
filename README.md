# db-exporter

db-exporter - app for export db schema to formats:
* CSV: export table data `csv`
* Markdown: export table structure `md`
* Class diagram: export table structure `diagram`
* Go structures with db tags `go-structs`
* Goose migrations `goose`
* Migrations for [sql-migrate](https://github.com/rubenv/sql-migrate) `go-sql-migrate`
* Raw SQL Laravel migrations `laravel-migrations-raw`
* Laravel models `laravel-models`
* YAML fixtures `yaml-fixtures`

Supported database: PostgreSQL

usage:
```text
Usage
  db-exporter[--config] [--tasks]

Options
  config                        Path to config file (yaml), default: ./.db-exporter.yaml
  tasks                         task names of config file

Usage examples
  db-exporter --config db.yaml  Run db-exporter with custom config path
```

Config file declared in [JSON Schema](db-exporter-json-schema.json)

**Export from postgres to markdown**

Add config file as `.db-exporter.yaml`

```yaml
databases:
  default:
    driver: postgres
    dsn: ${PG_DSN}

tasks:
  gen_md:
    activities:
      - export: md
        spec:
          with_diagram: true
        out:
          dir: ./out

options:
  print_stat: true
  debug: true
```

Run: `db-exporter`

**Export/import with YAML**

Add config file as `.db-exporter.yaml`
```yaml
databases:
  default:
    driver: postgres
    dsn: $PG_DSN

tasks:
  export:
    activities:
      - export: yaml-fixtures
        out:
          dir: ./data

  import:
    activities:
      - import: yaml-fixtures
        from: ./data
```

Run export: `db-exporter --tasks=export`

Run import: `db-exporter --tasks=import`

## Using custom templates

[Twig syntax](https://twig.symfony.com) is used to compile templates. The Twig port is a [Stick](https://github.com/tyler-sommer/stick).

| Exporter               | Template                     | Description                                                  |
|------------------------|------------------------------|--------------------------------------------------------------|
| csv                    | csv/export_single.csv        | Template for generate single csv file                        |
| diagram                | diagram/table.html           | Template for generate table                                  |
| go-sql-migrate         | go-sql-migrate/migration.sql | Template for generate migration                              |
| go-structs             | go-structs/model.go.tpl      | Template for generate table                                  |
| goose                  | goose/migration.sql          | Template for generate migration                              |
| goose-fixtures         | goose/migration.sql          | Template for generate migration with fixtures                |
| grpc-crud              | grpc-crud/gprc.proto         | Template for generate protobuf                               |
| laravel-migrations-raw | laravel/migration-raw.php    | Template for generate migration                              |
| laravel-models         | laravel/model.php            | Template for generate model                                  |
| md                     | md/single-tables.md          | Template for generate single markdown file                   |
| md                     | md/per-index.md              | Template for generate index markdown file (--table-per-file) |
| md                     | md/per-table.tmd             | Template for generate table markdown file (--table-per-file) |

You can download templates from [/templates](./templates)

In order for the db-exporter to use **your** templates, you need to place them in the `./db-exporter-templates` folder

## Use with GitHub Actions

You can run `db-exporter` as a GitHub action.

Add config file `.db-exporter.yaml`:
```yaml
databases:
  default:
    driver: postgres
    dsn: ${PG_DSN}

tasks:
  gen_docs:
    commit:
      message: "[auto] add documentation for database schema"
      push: true
    activities:
      - export: md
        spec:
          with_diagram: true
        out:
          dir: ./docs
```

Add GitHub Action as `./.github/workflows/docs.yaml`
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

      # step for run migration

      - name: Generate markdown docs
        uses: artarts36/db-exporter@master
        env:
          PG_DSN: "host=localhost port=5499 user=test password=test dbname=cars sslmode=disable"
        with:
          tasks: gen_docs
````
