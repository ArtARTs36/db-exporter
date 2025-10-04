# Usage examples

## Export schema from PostgreSQL to Markdown

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
```

Run: `PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter`

## Export/import data to YAML

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

Run export: `$PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter --tasks=export`

Run import: `$PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter --tasks=import`

## Export schema to Go entities and repositories

Add config file as `.db-exporter.yaml`
```yaml
databases:
  default:
    driver: postgres
    dsn: $PG_DSN

tasks:
  export:
    activities:
      - export: go-entity-repository
        skip_exists: true
        spec:
          entities:
            package: internal/domain
          repositories:
            package: internal/infrastructure/repositories
            container:
              struct_name: group
            interfaces:
              place: entity
            with_mocks: true
        out:
          dir: ./ # is root project path

  import:
    activities:
      - import: yaml-fixtures
        from: ./data
```

Run export: `$PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter --tasks=export`

Run import: `$PG_DSN="port=5459 user=db password=db dbname=db sslmode=disable" db-exporter --tasks=import`
