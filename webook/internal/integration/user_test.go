package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"geek-basic-go/webook/internal/integration/startup"
	"geek-basic-go/webook/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// init
// 包初始化方法，一旦要执行使用包里的任何方法或变量，就会先提前执行这个方法，只执行一次
// 线程安全
func init() {
	// 少输出日志
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandler_SendSmsCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name string
		// 准备数据
		before func(t *testing.T)
		//验证数据
		after      func(t *testing.T)
		phone      string
		wantedCode int
		wantedBody ginx.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 准备数据
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) == 6)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:      "13813382231",
			wantedCode: http.StatusOK,
			wantedBody: ginx.Result{
				Code: 0,
				Msg:  "短信发送成功",
				Data: nil,
			},
		},
		{
			name: "未输入手机号码",
			before: func(t *testing.T) {
				// 准备数据
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
			},
			phone:      "",
			wantedCode: http.StatusOK,
			wantedBody: ginx.Result{
				Code: 4,
				Msg:  "请输入正确手机号",
				Data: nil,
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 准备数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second+50).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				// 把数据拿出来后再删除掉
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:      "13813382231",
			wantedCode: http.StatusOK,
			wantedBody: ginx.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
				Data: nil,
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				// 准备数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				// redis中有key但是没有过期时间，web/user.go中的SendSmsLoginCode代码会进入default系统错误分支
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据
				// 需要验证验证码存在了redis里
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:13813382231"
				// 把数据拿出来后再删除掉
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:      "13813382231",
			wantedCode: http.StatusOK,
			wantedBody: ginx.Result{
				Code: 5,
				Msg:  "系统错误",
				Data: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Println(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))
			tc.before(t)
			defer tc.after(t)
			// 准备请求
			req, err := http.NewRequest(http.MethodPost, "/users/login/sms/code",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone":"%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 调用请求
			server.ServeHTTP(recorder, req)
			// Check the response code
			assert.Equal(t, tc.wantedCode, recorder.Code)
			// Check the response body
			if tc.wantedCode != http.StatusOK {
				return
			}
			var res ginx.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantedBody, res)
		})
	}
}
