package modifiers

import (
	"fmt"
	"github.com/artarts36/gds"
	"strings"

	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Comments struct {
}

func (c *Comments) ModifyProcedure(p *presentation.Procedure) {
	comment := c.generateCommentForProcedure(p)

	p.SetCommentTop(comment)
}

func (c *Comments) ModifyService(s *presentation.Service) {
	_, pluralName := c.generateEntityNames(s.TableMessage().Table())

	s.SetCommentTop(fmt.Sprintf("Service for working with %s.", pluralName))
}

func (c *Comments) generateCommentForProcedure(p *presentation.Procedure) string {
	singleName, pluralName := c.generateEntityNames(p.Service().TableMessage().Table())

	switch p.Type() {
	case presentation.ProcedureTypeCreate:
		return fmt.Sprintf("Create a new %s.", singleName)
	case presentation.ProcedureTypeList:
		return fmt.Sprintf("Get all %s.", pluralName)
	case presentation.ProcedureTypeGet:
		return fmt.Sprintf("Get a %s.", singleName)
	case presentation.ProcedureTypeDelete:
		return fmt.Sprintf("Delete a %s.", singleName)
	case presentation.ProcedureTypeUpdate:
		return fmt.Sprintf("Update an existing %s.", singleName)
	case presentation.ProcedureTypeUndelete:
		return fmt.Sprintf("Restore a soft-deleted %s.", singleName)
	default:
		return ""
	}
}

func (c *Comments) generateEntityNames(table *schema.Table) (singleName string, pluralName string) {
	words := table.Name.SplitWords()

	single := strings.Builder{}
	plural := strings.Builder{}

	for i, word := range words {
		w := strings.ToLower(word.Word)

		if i < len(words)-1 {
			single.WriteString(w + " ")
			plural.WriteString(w + " ")
		} else {
			single.WriteString(gds.NewString(w).Singular().Value)
			plural.WriteString(gds.NewString(w).Plural().Value)
		}
	}

	return single.String(), plural.String()
}
