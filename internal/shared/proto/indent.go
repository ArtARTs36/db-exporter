package proto

import "strings"

type Indent struct {
	step string
	curr string

	next *Indent
}

func NewIndent(step int) *Indent {
	return &Indent{
		step: strings.Repeat(" ", step),
		curr: "",
	}
}

func (i *Indent) Next() *Indent {
	if i.next == nil {
		i.next = &Indent{step: i.step, curr: i.curr + i.step}
	}

	return i.next
}
