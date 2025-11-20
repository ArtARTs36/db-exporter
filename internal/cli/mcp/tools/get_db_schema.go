package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/cli/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/schema"
	schema2 "github.com/artarts36/db-exporter/internal/schema"

	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

type GetDBSchemaTool struct {
	cfg  *config.Config
	info protocol.ToolInfo
}

func NewGetDBSchemaTool(cfg *config.Config) *GetDBSchemaTool {
	return &GetDBSchemaTool{
		cfg: cfg,
		info: protocol.ToolInfo{
			Name:        "get-database-schema",
			Description: "get schema of database",
			InputSchema: jsonschema.Property{
				Type: jsonschema.TypeObject,
				Properties: map[string]jsonschema.Property{
					"database_name": {
						Title: "Name of database from db-exporter configuration file",
						Type:  jsonschema.TypeString,
					},
				},
			},
		},
	}
}

func (t *GetDBSchemaTool) Info() protocol.ToolInfo {
	return t.info
}

func (t *GetDBSchemaTool) Execute(ctx context.Context, args map[string]interface{}) (any, error) {
	db, err := t.selectDatabase(args)
	if err != nil {
		return nil, err
	}

	con := conn.NewConnection(db)

	sch, err := schema.Load(ctx, con)
	if err != nil {
		return nil, err
	}

	return json.Marshal(sch)
}

func (t *GetDBSchemaTool) selectDatabase(args map[string]interface{}) (schema2.Database, error) {
	dbName, ok := args["database_name"].(string)
	if ok {
		if db, dbOk := t.cfg.Databases[dbName]; dbOk {
			return db, nil
		}
		return schema2.Database{}, fmt.Errorf("database %q is not defined", dbName)
	}

	db, ok := t.cfg.GetDefaultDatabase()
	if ok {
		return db, nil
	}

	return schema2.Database{}, errors.New("databases not defined")
}
