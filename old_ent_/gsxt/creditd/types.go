package creditd

import (
	"gsxt/credit"
)

// 支持下发配置
type Request struct {
	Extra       string
	Name        string
	Keyword     string
	SearchCount int
	DetailCount int
}

// 支持返回多个结果
type Response struct {
	Extra   string          `json:"extra"`
	Name    string          `json:"name"`
	Keyword string          `json:keyword`
	State   string          `json:"state"`
	Keys    []string        `json:"keys"`
	Infos   []credit.InfoV2 `json:"infos"`
}
