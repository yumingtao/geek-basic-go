package web

import (
	"bytes"
	"context"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	svcmocks "geek-basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		// UserHandler用到了 UserService和CodeService，需要准备这两个mock实例
		// 定义了一个mock方法，返回这两个Service实例
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 构造预期请求输入
		reqBuilder func(t *testing.T) *http.Request
		// 预期response code
		wantedCode int
		// 预期response body
		wantedBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"ConfirmPassword": "hello#world123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "Hello, 恭喜注册成功",
		},
		{
			name: "Bind出错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 构造了一个非法json
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world12
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusBadRequest,
			wantedBody: "",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 构造了一个非法json
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123qq.com",
"password": "hello#world123",
"ConfirmPassword": "hello#world123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "邮箱格式不正确",
		},
		{
			name: "两次输入密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 构造了一个非法json
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"ConfirmPassword": "hello#world123aaa"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "两次输入密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 构造了一个非法json
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "helloworld123",
"ConfirmPassword": "helloworld123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "密码必须包含数字、特殊字符，并且长度不能小于8位",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("DB连接错误"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"ConfirmPassword": "hello#world123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "系统错误！",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrDuplicateEmail)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"ConfirmPassword": "hello#world123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantedCode: http.StatusOK,
			wantedBody: "注册用户失败:邮箱已被注册，请换个邮箱重新申请",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建ctrl
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// mock UserService 和 CodService
			userSvc, codeSvc := tc.mock(ctrl)
			// 创建UserHandler
			hdl := NewUserHandler(userSvc, codeSvc)
			// 注册路由
			server := gin.Default()
			hdl.RegisterRoutes(server)
			// 准备请求
			req := tc.reqBuilder(t)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 调用请求
			server.ServeHTTP(recorder, req)
			// Check the response code
			assert.Equal(t, tc.wantedCode, recorder.Code)
			// Check the response body
			assert.Equal(t, tc.wantedBody, recorder.Body.String())
		})
	}
}

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

func TestHttp(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("Login请求体")))
	t.Log(req)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestMock(t *testing.T) {
	// 先创建一个mock控制器
	ctrl := gomock.NewController(t)
	// 每个测试结束都要调用Finish，然后mock会验证测试流程是否符合预期
	defer ctrl.Finish()
	// mock 实现
	userSvc := svcmocks.NewMockUserService(ctrl)
	// 模拟调用Signup
	// 注意：涉及了几个模拟调用，在使用的时候就要都用上，而且顺序也要对上
	userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(nil)

	err := userSvc.SignUp(context.Background(), domain.User{Id: 1, Email: "123@qq.com"})
	assert.NoError(t, err)
}
