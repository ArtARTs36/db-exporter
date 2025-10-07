# Examples of using a `custom` exporter

## Generating .txt files from a template built into the configuration

Add config file as `.db-exporter.yaml`

```yaml
databases:
  default:
    driver: postgres
    dsn: $PG_DSN

tasks:
  gen_txt:
    activities:
      - format: custom
        spec:
          template: '{% for table in schema.Tables %}{{ table.Name.Value }},{% endfor %}'
          output:
            extension: "txt"
        out:
          dir: ./data
```

Run export: `PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter`

## Generating .txt files from a template from local file system

Add config file as `.db-exporter.yaml`

```yaml
databases:
  default:
    driver: postgres
    dsn: $PG_DSN

tasks:
  gen_txt:
    activities:
      - format: custom
        spec:
          template: '@local/path/to/file.txt'
          output:
            extension: "txt"
        out:
          dir: ./data
```

Run export: `PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter`
