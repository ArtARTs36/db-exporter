package golang

import (
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"slices"
)

type ImportGroup struct {
	*ds.Set[string]
}

type ImportGroups struct {
	std    *ImportGroup
	shared *ImportGroup
	local  *ImportGroup
}

func NewImportGroups() *ImportGroups {
	return &ImportGroups{
		std:    &ImportGroup{Set: ds.NewSet[string]()},
		shared: &ImportGroup{Set: ds.NewSet[string]()},
		local:  &ImportGroup{Set: ds.NewSet[string]()},
	}
}

func (g *ImportGroup) SortedList() []string {
	list := g.List()
	slices.Sort(list)
	return list
}

func (g *ImportGroups) AddStd(pkg string) {
	g.std.Add(pkg)
}

func (g *ImportGroups) AddShared(pkg string) {
	g.shared.Add(pkg)
}

func (g *ImportGroups) AddLocal(pkg string) {
	g.local.Add(pkg)
}

func (g *ImportGroups) Sorted() [][]string {
	const groupCapacity = 3

	groups := make([][]string, 0, groupCapacity)

	if g.std.Valid() {
		groups = append(groups, g.std.SortedList())
	}

	if g.shared.Valid() {
		groups = append(groups, g.shared.SortedList())
	}

	if g.local.Valid() {
		groups = append(groups, g.local.SortedList())
	}

	return groups
}

func (g *ImportGroups) Valid() bool {
	return g.std.Valid() || g.shared.Valid() || g.local.Valid()
}
