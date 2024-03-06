# db-exporter

db-exporter - simple app for export db schema to formats:
* markdown `md`
* class diagram `diagram`
* go structs `go-structs`
* goose migrations `goose`
* raw laravel migrations `laravel-migrations-raw`

usage:
```text
./db-exporter driver-name dsn format out-dir [--table-per-file] [--with-diagram]

Arguments
  driver-name                database driver name, required, available values: [pg, postgres]
  dsn                        data source name, required
  format                     exporting format, required, available values: [md, diagram, go-structs, goose, laravel-migrations-raw]
  out-dir                    Output directory, required

Options
  table-per-file             Export one table to one file
  with-diagram               Export with diagram (only markdown)
  without-migrations-table   Export without migrations table
```

**Export from postgres to markdown**

```db-exporter pg "host=postgres user=root password=root dbname=cars" md ./docs```