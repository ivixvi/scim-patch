package scimpatch_test

import (
	"fmt"

	"github.com/scim2/filter-parser/v2"
)

// path は *filter.Path を取得しやすくするためのテストユーティリティです。
// APIのリクエストボディ Operations[].path にはいってくる想定の値を引数に与えて利用します。
func path(s string) *filter.Path {
	p, err := filter.ParsePath([]byte(s))
	if err != nil {
		fmt.Printf("Failed to parse %s occurred by %s\n", s, err)
	}
	return &p
}
