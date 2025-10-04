# db-exporter

![Docker Image Version](https://img.shields.io/docker/v/artarts36/db-exporter?style=for-the-badge&logo=docker&label=Image%20Version&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fartarts36%2Fdb-exporter)
![Docker Image Size](https://img.shields.io/docker/image-size/artarts36/db-exporter?style=for-the-badge&logo=docker&label=Image%20Size&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fartarts36%2Fdb-exporter)

**db-exporter** - application for export db schema and data to formats:
* CSV: export table data `csv`
* Markdown: export table structure `md`
* Class diagram: export table structure `diagram`
* Go structures with db tags `go-entities`
* Goose migrations `goose`
* Goose Fixtures: Goose migrations with inserts `goose-fixtures`
* Migrations for [sql-migrate](https://github.com/rubenv/sql-migrate) `go-sql-migrate`
* Raw SQL Laravel migrations `laravel-migrations-raw`
* Laravel models `laravel-models`
* YAML fixtures `yaml-fixtures`
* Go Entity Repository `go-entity-repository`
* JSON Schema `json-schema`
* GraphQL `graphql`
* DBML `dbml`: export db schema (table, ref, enum) to dbml

Supported database schemas:
- PostgreSQL: full-support
- DBML: only export db schema, without fixtures and import
- MySQL: only generate migrations from PostgreSQL/DBML

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

▷ Usage examples:
- [🚀Use with GitHub Actions](./docs/usage_examples.md#use-with-github-actions)
- [Export schema from PostgreSQL to Markdown](./docs/usage_examples.md#export-schema-from-postgresql-to-markdown)
- [Export/import data to YAML](./docs/usage_examples.md#exportimport-data-to-yaml)
- [Export schema to Go entities and repositories](./docs/usage_examples.md#export-schema-to-go-entities-and-repositories)

## Environment variables
You can inject environment variables to config:

- **DSN** to database in `databases`:
    ```yaml
    databases:
      default:
        driver: postgres
        dsn: $PG_DSN
    ```
- **Commit Author** in `tasks`:
    ```yaml
  tasks:
    gen_docs:
      commit: 
        author: ${COMMIT_AUTHOR}
      activities:
        ...
    ```
- **Tables List** in `activities`
    ```yaml
  tasks:
    gen_csv:
      activities:
        - export: csv
          tables:
            list: ${MY_TABLES}
    ```
