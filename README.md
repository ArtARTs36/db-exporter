# db-exporter

db-exporter - simple app for export db schema to formats:
* markdown
* class diagram

usage:
```text
./db-exporter driver-name dsn format out-dir [--table-per-file] [--with-diagram]

Arguments
  driver-name       database driver name, required, available values: [pg, postgres]
  dsn               data source name, required
  format            exporting format, required, available values: [md, diagram]
  out-dir           Output directory, required

Options
  table-per-file    Export one table to one file
  with-diagram      Export with diagram (only markdown)
```
