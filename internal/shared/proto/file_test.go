package proto

import (
	"github.com/artarts36/gds"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Render(t *testing.T) {
	f := &File{
		Package: "my-super-package",
		Services: []*Service{
			{
				Name: "UserService",
				Procedures: []*ServiceProcedure{
					{
						Name:    "Get",
						Param:   "GetUserRequest",
						Returns: "GetUserResponse",
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
			NewEnumWithValues(gds.NewString("UserStatus"), []string{
				"ACTIVE",
				"BANNED",
			}),
		},
	}

	expected := `syntax = "proto3";

package my-super-package;

service UserService {
    rpc Get(GetUserRequest) returns (GetUserResponse) {}
}

message GetUserRequest {
    int64 id = 1;
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
	got := f.Render()

	assert.Equal(t, expected, got)
}
