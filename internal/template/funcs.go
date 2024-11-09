package template

import (
	"fmt"
	"github.com/artarts36/gds"
	"strings"
	"time"

	"github.com/tyler-sommer/stick"
)

func twigFuncs() map[string]stick.Func {
	const spacesAfterArgsCount = 2

	noCtx := func(fn func(args ...stick.Value) stick.Value) stick.Func {
		return func(_ stick.Context, args ...stick.Value) stick.Value {
			return fn(args...)
		}
	}

	return map[string]stick.Func{
		"bool_string": noCtx(func(args ...stick.Value) stick.Value {
			val, _ := args[0].(bool)
			if val {
				return "true"
			}
			return "false"
		}),
		"quote_string": noCtx(func(args ...stick.Value) stick.Value {
			switch val := args[0].(type) {
			case string:
				return fmt.Sprintf("%q", val)
			case time.Time:
				if val.IsZero() {
					return ""
				}

				return fmt.Sprintf("%q", val.Format(time.RFC3339))
			}
			return args[0]
		}),
		"spaces_after": noCtx(func(args ...stick.Value) stick.Value {
			if len(args) != spacesAfterArgsCount {
				return ""
			}

			var currentStringLen int

			switch v := args[0].(type) {
			case string:
				currentStringLen = len(v)
			case *gds.String:
				currentStringLen = v.Len()
			default:
				return ""
			}

			needLength, valid := args[1].(int)
			if !valid {
				return ""
			}

			repeats := needLength - currentStringLen
			if repeats == 0 {
				return ""
			}

			return strings.Repeat(" ", repeats)
		}),
	}
}
