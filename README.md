# db-exporter

db-exporter - simple app for export db schema to formats:
* markdown

usage:
```text
./db-exporter driver-name dsn format out-dir [--table-per-file]

Arguments
  driver-name                                                                database driver name, required, available values: [pg, postgres]
  dsn                                                                        data source name, required
  format                                                                     exporting format, required, available values: [md]
  out-dir                                                                    Output directory, required

Options
  table-per-file                                                             Export one table to one file
```
