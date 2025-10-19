package proto

import (
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"github.com/artarts36/gds"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Render(t *testing.T) {
	f := &File{
		Package: "my-super-package",
		Imports: gds.NewSet[string]("google/protobuf/timestamp.proto"),
		Services: []*Service{
			{
				Name: "UserService",
				Procedures: []*ServiceProcedure{
					{
						Name:    "List",
						Param:   "GetUserRequest",
						Returns: "GetUserResponse",
						Options: []*ServiceProcedureOption{
							{
								Name: "google.api.http",
								Params: map[string]string{
									"get": "/v1/users",
								},
							},
						},
					},
					{
						Name:    "Get",
						Param:   "GetUserRequest",
						Returns: "GetUserResponse",
						Options: []*ServiceProcedureOption{
							{
								Name: "google.api.http",
								Params: map[string]string{
									"get": "/v1/users/{id}",
								},
							},
						},
					},
				},
			},
		},
		Messages: []*Message{
			{
				Name: "GetUserRequest",
				Fields: []*Field{
					{
						Type: "int64",
						Name: "id",
						ID:   1,
						Options: []*FieldOption{
							{
								Name:  "google.api.field_behavior",
								Value: "REQUIRED",
							},
						},
					},
				},
			},
			{
				Name: "GetUserResponse",
				Fields: []*Field{
					{
						Type: "int64",
						Name: "id",
						ID:   1,
					},
					{
						Type: "string",
						Name: "name",
						ID:   2,
					},
				},
			},
		},
		Enums: []*Enum{
			NewEnumWithValues(*gds.NewString("UserStatus"), []string{
				"ACTIVE",
				"BANNED",
			}),
		},
	}

	expected := `syntax = "proto3";

package my-super-package;

import "google/protobuf/timestamp.proto";

service UserService {
  rpc List(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = {
      get: "/v1/users"
    };
  }

  rpc Get(GetUserRequest) returns (GetUserResponse) {
    option (google.api.http) = {
      get: "/v1/users/{id}"
    };
  }
}

message GetUserRequest {
  int64 id = 1 [(google.api.field_behavior) = REQUIRED];
}

message GetUserResponse {
  int64 id = 1;
  string name = 2;
}

enum UserStatus {
  USERSTATUS_UNDEFINED = 0;
  USERSTATUS_ACTIVE = 1;
  USERSTATUS_BANNED = 2;
}`
	got := f.Render(indentx.NewIndent(2))

	assert.Equal(t, expected, got)
}
