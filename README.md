# db-exporter

db-exporter - simple app for export db schema to formats:
* Markdown `md`
* Class diagram `diagram`
* Go structures with db tags `go-structs`
* Goose migrations `goose`
* Migrations for golang-migrate `golang-migrate`
* Raw SQL Laravel migrations `laravel-migrations-raw`

usage:
```text
./db-exporter driver-name dsn format out-dir [--table-per-file] [--with-diagram] [--without-migrations-table] [--tables=<table_name>,...] [--package=<package name>]

Arguments
  driver-name                database driver name, required, available values: [pg, postgres]
  dsn                        data source name, required
  format                     exporting format, required, available values: [md, diagram, go-structs, goose, laravel-migrations-raw]
  out-dir                    Output directory, required

Options
  table-per-file             Export one table to one file
  with-diagram               Export with diagram (only markdown)
  without-migrations-table   Export without migrations table
  tables                     Table list for export, separator: ","
  package                    Package name for code gen, e.g: models
```

**Export from postgres to markdown**

```db-exporter pg "host=postgres user=root password=root dbname=cars" md ./docs```
