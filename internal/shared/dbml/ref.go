package dbml

import "fmt"

type Ref struct {
	From string
	Type string
	To   string
}

func (r *Ref) Render() string {
	return fmt.Sprintf("Ref: %s %s %s", r.From, r.Type, r.To)
}
