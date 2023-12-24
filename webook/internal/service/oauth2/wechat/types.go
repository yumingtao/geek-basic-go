package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"geek-basic-go/webook/internal/domain"
	"net/http"
	"net/url"
)

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

const (
	authURLPattern        = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	accessTokenUrlPattern = `https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`
)

type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type WechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func (w *WechatService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(accessTokenUrlPattern, w.appId, w.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpRes, err := w.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(httpRes.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败，errcode %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}

	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

func NewService(appId string, appSecret string) Service {
	return &WechatService{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (w *WechatService) AuthUrl(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, w.appId, redirectURL, state), nil
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires-in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openId"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionId"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
