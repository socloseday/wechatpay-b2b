package example

import (
	"context"

	"github.com/enjoy322/wechatpay-b2b/client"
	"github.com/enjoy322/wechatpay-b2b/service"
	"github.com/enjoy322/wechatpay-b2b/types"
)

// GetBalanceExample 展示如何查询商户号账户余额（微信支付手动提现流程中的“查询账户余额”接口）。
//
// 外部只需定时更新 c.AccessToken / c.AppKey 即可（本函数不需要每次传入 access_token 和 appKey）。
func GetBalanceExample(ctx context.Context, c *client.Client, mchid string) (*types.BalanceResponse, error) {
	svc := service.NewBalanceService(c)
	return svc.GetBalance(ctx, types.BalanceRequest{Mchid: mchid})
}
