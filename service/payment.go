package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/enjoy322/wechatpay-b2b/client"
	"github.com/enjoy322/wechatpay-b2b/types"
)

// PaymentService 负责构建小程序支付相关参数。
type PaymentService interface {
	// BuildCombinedPaymentParams 生成合单支付参数，用于小程序 wx.requestCommonPayment。
	BuildCombinedPaymentParams(ctx context.Context, req types.CombinedPaymentSignData) (*types.CommonPaymentParams, error)
}

type paymentService struct {
	client *client.Client
}

const combinePayURI = "/retail/B2b/combinepay"

// NewPaymentService 创建支付服务。
func NewPaymentService(c *client.Client) PaymentService {
	return &paymentService{client: c}
}

// BuildCombinedPaymentParams 生成合单支付参数，用于小程序 wx.requestCommonPayment。
func (s *paymentService) BuildCombinedPaymentParams(ctx context.Context, req types.CombinedPaymentSignData) (*types.CommonPaymentParams, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if s.client.GetAppKey() == "" {
		return nil, errors.New("appKey is empty")
	}
	if len(req.CombinedOrderList) == 0 {
		return nil, errors.New("combined_order_list is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return &types.CommonPaymentParams{
		SignData:  string(body),
		Mode:      "0",
		PaySig:    s.client.GetPaySig(combinePayURI, body),
		Signature: s.client.GetUserSignature(body),
	}, nil
}
