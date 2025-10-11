package pg

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/pg/tokens"
	"github.com/bzick/tokenizer"
)

type Query interface {
	query()
}

type CreateTableQuery struct {
	Table string

	Columns []*Column
}

type Column struct {
	Name     string
	DataType string

	IsPrimaryKey bool
	Nullable     bool
}

func (q *CreateTableQuery) query() {}

func (q *CreateTableQuery) parse(stream *tokenizer.Stream) error {
	stream.GoNext()

	if stream.CurrentToken().Is(tokenizer.TokenKeyword) {
		q.Table = stream.CurrentToken().ValueString()
	} else {
		return unexpectedToken(stream.CurrentToken())
	}

	stream.GoNext()
	if !stream.CurrentToken().Is(tokens.BracketLeft) {
		return unexpectedToken(stream.CurrentToken())
	}

	for stream.IsValid() && !stream.NextToken().Is(tokens.BracketRight) {
		stream.GoNext()

		col := &Column{}

		if err := col.parse(stream); err != nil {
			return fmt.Errorf("parse column: %w", err)
		}

		q.Columns = append(q.Columns, col)
	}

	return nil
}

func (c *Column) parse(stream *tokenizer.Stream) error {
	for stream.IsValid() {
		if stream.CurrentToken().Is(tokenizer.TokenKeyword) {
			if c.Name == "" {
				c.Name = stream.CurrentToken().ValueString()
			} else {
				if c.DataType == "" {
					c.DataType = stream.CurrentToken().ValueString()
				} else {
					c.DataType += " " + stream.CurrentToken().ValueString()
				}
			}
		}

		if stream.CurrentToken().Is(tokens.PrimaryKey) {
			c.IsPrimaryKey = true
		}

		if stream.NextToken().Is(tokens.Comma, tokens.BracketRight) {
			return nil
		}

		stream.GoNext()
	}

	return nil
}
