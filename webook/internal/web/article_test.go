package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"geek-basic-go/webook/internal/domain"
	"geek-basic-go/webook/internal/service"
	svcmocks "geek-basic-go/webook/internal/service/mocks"
	ijwt "geek-basic-go/webook/internal/web/jwt"
	"geek-basic-go/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		reqBody string

		wantedCode int
		wantedRes  Result
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody:    `{"title": "我的标题","content": "我的内容"}`,
			wantedCode: 200,
			wantedRes: Result{
				Data: float64(1),
			},
		},
		{
			name: "已经帖子发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      123,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(123), nil)
				return svc
			},
			reqBody:    `{"id": 123, "title": "我的标题","content": "我的内容"}`,
			wantedCode: 200,
			wantedRes: Result{
				Data: float64(123),
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("发表失败"))
				return svc
			},
			reqBody:    `{"title": "我的标题","content": "我的内容"}`,
			wantedCode: 200,
			wantedRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				return svc
			},
			// 注意此处title前边少了一个"
			reqBody:    `{title": "我的标题","content": "我的内容"}`,
			wantedCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// mock UserService 和 CodService
			svc := tc.mock(ctrl)
			// 创建UserHandler
			hdl := NewArticleHandler(svc, logger.NewNopLogger())
			// 注册路由
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UserClaims{
					Uid: 123,
				})
			})
			hdl.RegisterRoutes(server)
			// 准备请求
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBufferString(tc.reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 调用请求
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantedCode, recorder.Code)
			if tc.wantedCode != http.StatusOK {
				return
			}
			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			// Check the response code
			// Check the response body
			assert.Equal(t, tc.wantedRes, res)
		})
	}
}
