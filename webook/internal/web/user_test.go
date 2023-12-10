package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserEmailPattern(t *testing.T) {
	// Table Driven
	testCase := []struct {
		// 测试用例结构定义
		name  string
		email string
		match bool
	}{ // 测试用例实例
		{
			name:  "不带@",
			email: "123456_126.com",
			match: false,
		},
		{
			name:  "带@但没有后后缀",
			email: "123456@126",
			match: false,
		},
		{
			name:  "合法邮箱",
			email: "123456@126.com",
			match: true,
		},
	}
	h := NewUserHandler(nil, nil)
	// 执行测试用例
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRexExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}
