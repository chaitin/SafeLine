package model

type Behaviour struct {
	Base
	SrcRouter string `json:"src_router"`
	DstRouter string `json:"dst_router"`
}
