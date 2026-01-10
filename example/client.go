package example

import (
	"github.com/enjoy322/wechatpay-b2b/client"
)

// NewClientExample 展示如何初始化 SDK Client（用于 service.NewBalanceService / service.NewRetailService 等）。
// 说明：
// - AccessToken / AppKey 都保存在 client 内，业务侧只需定时刷新并更新到同一个 client 上。
func NewClientExample(accessToken, appKey string) (*client.Client, error) {
	return client.NewClient(client.Options{
		BaseURL:        "https://api.weixin.qq.com",
		AccessToken:  accessToken,
		AppKey: appKey,
	})
}
