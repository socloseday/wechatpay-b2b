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
	// BuildPaymentParams 生成单订单支付参数，用于小程序 wx.requestCommonPayment。
	BuildPaymentParams(ctx context.Context, req types.Order, sessionKey string) (*types.CommonPaymentParams, error)
	// BuildCombinedPaymentParams 生成合单支付参数，用于小程序 wx.requestCommonPayment。
	BuildCombinedPaymentParams(ctx context.Context, req types.CombinedPaymentSignData, sessionKey string) (*types.CommonPaymentParams, error)
}

type paymentService struct {
	client *client.Client
}

const (
	requestCommonPaymentURI = "requestCommonPayment"
	paymentModeGoods        = "retail_pay_goods"
	paymentModeCombined     = "retail_pay_combined_goods"
)

// NewPaymentService 创建支付服务。
func NewPaymentService(c *client.Client) PaymentService {
	return &paymentService{client: c}
}

// BuildPaymentParams 生成单订单支付参数，用于小程序 wx.requestCommonPayment。
func (s *paymentService) BuildPaymentParams(ctx context.Context, req types.Order, sessionKey string) (*types.CommonPaymentParams, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if s.client.GetAppKey() == "" {
		return nil, errors.New("appKey is empty")
	}
	if sessionKey == "" {
		return nil, errors.New("sessionKey is empty")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return &types.CommonPaymentParams{
		SignData:  string(body),
		Mode:      paymentModeGoods,
		PaySig:    s.client.GetPaySig(requestCommonPaymentURI, body),
		Signature: s.client.GetUserSignature(body, sessionKey),
	}, nil
}

// BuildCombinedPaymentParams 生成合单支付参数，用于小程序 wx.requestCommonPayment。
func (s *paymentService) BuildCombinedPaymentParams(ctx context.Context, req types.CombinedPaymentSignData, sessionKey string) (*types.CommonPaymentParams, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if s.client.GetAppKey() == "" {
		return nil, errors.New("appKey is empty")
	}
	if sessionKey == "" {
		return nil, errors.New("sessionKey is empty")
	}
	if len(req.CombinedOrderList) == 0 {
		return nil, errors.New("combined_order_list is required")
	}

	type paySigItem struct {
		Mchid  string `json:"mchid"`
		PaySig string `json:"paysig"`
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paySigItems := make([]paySigItem, 0, len(req.CombinedOrderList))
	for _, order := range req.CombinedOrderList {

		paySigPer := s.client.GetPaySig(requestCommonPaymentURI, body)

		paySigItems = append(paySigItems, paySigItem{
			Mchid:  order.Mchid,
			PaySig: paySigPer,
		})
	}

	paySigBytes, err := json.Marshal(paySigItems)
	if err != nil {
		return nil, err
	}

	return &types.CommonPaymentParams{
		SignData:  string(body),
		Mode:      paymentModeCombined,
		PaySig:    string(paySigBytes),
		Signature: s.client.GetUserSignature(body, sessionKey),
	}, nil
}
