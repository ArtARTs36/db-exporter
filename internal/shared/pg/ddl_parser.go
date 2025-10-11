package pg

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/shared/pg/tokens"
	"github.com/bzick/tokenizer"
)

type DDLParser struct {
	tokenizer *tokenizer.Tokenizer
}

type DDL struct {
	Queries []Query
}

func NewDDLParser() *DDLParser {
	tkn := tokenizer.New()
	tkn.DefineTokens(tokens.CreateTable, []string{"CREATE TABLE"})
	tkn.DefineTokens(tokens.BracketLeft, []string{"("})
	tkn.DefineTokens(tokens.BracketRight, []string{")"})
	tkn.DefineTokens(tokens.Comma, []string{","})
	tkn.DefineTokens(tokens.NotNull, []string{"NOT NULL", "not null", "NOT null", "not NULL"})
	tkn.DefineTokens(tokens.PrimaryKey, []string{"PRIMARY KEY", "primary key", "PRIMARY key", "primary KEY"})
	tkn.DefineTokens(tokens.Semicolon, []string{";"})
	tkn.AllowKeywordSymbols(tokenizer.Underscore, tokenizer.Numbers)

	return &DDLParser{
		tokenizer: tkn,
	}
}

func (p *DDLParser) Parse(ddlQuery string) (*DDL, error) {
	stream := p.tokenizer.ParseString(ddlQuery)
	defer stream.Close()

	result := &DDL{
		Queries: make([]Query, 0),
	}

	for stream.IsValid() {
		if stream.CurrentToken().Is(tokens.CreateTable) {
			tblQuery := &CreateTableQuery{}

			err := tblQuery.parse(stream)
			if err != nil {
				return nil, fmt.Errorf("parse create table query: %w", err)
			}

			result.Queries = append(result.Queries, tblQuery)
		}

		stream.GoNext()
	}

	return result, nil
}

func unexpectedToken(tkn *tokenizer.Token) error {
	return fmt.Errorf("unexpected token: %v", tkn)
}
