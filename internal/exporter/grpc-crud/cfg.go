package grpccrud

import (
	"github.com/artarts36/specw"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type paginationType string

const (
	paginationTypeOffset paginationType = "offset"
	paginationTypeToken  paginationType = "token"
	paginationTypeNone   paginationType = "none"
)

type Specification struct {
	Package    string                                     `yaml:"package" json:"package"`
	Indent     int                                        `yaml:"indent" json:"indent"`
	Options    orderedmap.OrderedMap[string, interface{}] `yaml:"options" json:"options"`
	Pagination paginationType                             `yaml:"pagination" json:"pagination"`
	With       specw.BoolObject[struct {
		GoogleApiHttp specw.BoolObject[struct { //nolint:revive // <- not readable
			PathPrefix string `yaml:"path_prefix" json:"path_prefix"`
		}] `yaml:"google.api.http" json:"google.api.http"`
		GoogleAPIFieldBehavior specw.BoolObject[struct{}] `yaml:"google.api.field_behavior" json:"google.api.field_behavior"`
		BufValidateField       specw.BoolObject[struct{}] `yaml:"buf.validate.field" json:"buf.validate.field"`
	}] `yaml:"with" json:"with"`
}

func (s *Specification) Validate() error {
	if s.Indent == 0 {
		s.Indent = 2
	}

	if s.Pagination == "" {
		s.Pagination = paginationTypeOffset
	}

	return nil
}
