package fsl

import (
	"fmt"
)

type Selector struct {
	TableName string
	ID        string
	Where     string
	Action    string
}

func NewSelector(tableName, id, where, action string) *Selector {
	return &Selector{
		TableName: tableName,
		ID:        id,
		Where:     where,
		Action:    action,
	}
}

func (s *Selector) String() string {
	ret := fmt.Sprintf("ID %s", Quote(s.ID))
	if s.Where != "" {
		ret += " WHERE " + s.Where
	}
	ret += fmt.Sprintf(" ACTION %s", s.Action)
	return ret
}

func (s *Selector) AppendInto() string {
	return fmt.Sprintf("APPEND INTO %s %s;", s.TableName, s.String())
}

func AppendInto(tableName, id, where, action string) string {
	return NewSelector(tableName, id, where, action).AppendInto()
}
