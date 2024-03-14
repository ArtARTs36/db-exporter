# db-exporter

db-exporter - simple app for export db schema to formats:
* Markdown `md`
* Class diagram `diagram`
* Go structures with db tags `go-structs`
* Goose migrations `goose`
* Migrations for [sql-migrate](https://github.com/rubenv/sql-migrate) `go-sql-migrate`
* Raw SQL Laravel migrations `laravel-migrations-raw`

Supported database: PostgreSQL

usage:
```text
db-exporter driver-name dsn format out-dir [--table-per-file] [--with-diagram] [--without-migrations-table] [--tables=<value>] [--package=<value>] [--file-prefix=<value>]
[--commit-message=<value>] [--commit-push] [--commit-author=<value>] [--stat] [--debug]

Arguments
  driver-name                database driver name, required, available values: [pg]
  dsn                        data source name, required
  format                     exporting format, required, available values: [md, diagram, go-structs, goose, go-sql-migrate, laravel-migrations-raw]
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
| go-sql-migrate | go-sql-migrate/migration.sql | Template for generate migration                              |
| laravel        | laravel/migration-raw.php    | Template for generate migration                              |

You can download templates from [/templates](./templates)

In order for the db-exporter to use **your** templates, you need to place them in the `./db-exporter-templates` folder
