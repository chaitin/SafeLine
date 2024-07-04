package fsl

import (
	"fmt"
)

type Target struct {
	ID   string
	Type string
	Args string
}

func NewTarget(id string, typ string, args string) *Target {
	return &Target{
		ID:   id,
		Type: typ,
		Args: args,
	}
}

func NewSkynetTarget(id, configJSONStr string) *Target {
	return NewTarget(id, "skynet", configJSONStr)
}

func (t *Target) String() string {
	return fmt.Sprintf("%s TYPE %s ARGS (%s)", t.ID, t.Type, t.Args)
}

func (t *Target) Create() string {
	return fmt.Sprintf("CREATE TARGET %s;", t.String())
}

func CreateTarget(id, configJSONStr string) string {
	return NewSkynetTarget(id, configJSONStr).Create()
}
