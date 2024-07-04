package fsl

import (
	"fmt"
	"strings"
)

func Goto(m string) string {
	return fmt.Sprintf("GOTO %s", m)
}

func Wheres(items ...string) string {
	return strings.Join(items, " AND ")
}

func Actions(items ...string) string {
	return strings.Join(items, ", ")
}

func Set(k, t, v string) string {
	return fmt.Sprintf("SET %s::%s = %s", k, t, v)
}
