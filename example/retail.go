package example

import (
	"context"

	"github.com/enjoy322/wechatpay-b2b/client"
	"github.com/enjoy322/wechatpay-b2b/service"
	"github.com/enjoy322/wechatpay-b2b/types"
)

// BatchCreateRetailExample 展示如何预录入门店信息（/wxa/business/batchcreateretail）。
//
// 外部只需定时更新 c.AccessToken 即可（本函数不需要每次传入 access_token）。
func BatchCreateRetailExample(ctx context.Context, c *client.Client, list []types.RetailInfo) (*types.BatchCreateRetailResponse, error) {
	svc := service.NewRetailService(c)
	return svc.BatchCreateRetail(ctx, types.BatchCreateRetailRequest{RetailInfoList: list})
}
