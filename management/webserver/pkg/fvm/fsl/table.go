package fsl

import (
	"fmt"
	"strings"
)

type Table struct {
	Name     string
	Mappings map[State]State
}

func NewTable(name string, mappings map[State]State) *Table {
	return &Table{
		Name:     name,
		Mappings: mappings,
	}
}

func (t *Table) String() string {
	ret := fmt.Sprintf("TABLE %s", t.Name)
	if len(t.Mappings) > 0 {
		var mappings []string
		for from, to := range t.Mappings {
			mappings = append(mappings, fmt.Sprintf("%s TO %s", from.String(), to.String()))
		}
		ret += fmt.Sprintf(" MAPPING %s", strings.Join(mappings, ", "))
	}
	return ret
}

func (t *Table) Create() string {
	return fmt.Sprintf("CREATE %s;", t.String())
}

func CreateTable(name string, mappings map[State]State) string {
	return NewTable(name, mappings).Create()
}
