package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/enjoy322/wechatpay-b2b/client"
	"github.com/enjoy322/wechatpay-b2b/types"
)

// OrderService 处理订单查询、关单相关调用。
type OrderService interface {
	// CloseOrder 关闭订单。
	CloseOrder(ctx context.Context, req types.CloseOrderRequest) (*types.CloseOrderResponse, error)
	// GetOrder 查询订单。
	GetOrder(ctx context.Context, req types.GetOrderRequest) (*types.GetOrderResponse, error)
}

type orderService struct {
	client *client.Client
}

const (
	closeOrderURI = "/retail/B2b/closeb2border"
	getOrderURI   = "/retail/B2b/getorder"
)

// NewOrderService 创建订单服务。
func NewOrderService(c *client.Client) OrderService {
	return &orderService{client: c}
}

// CloseOrder 关闭订单。
func (s *orderService) CloseOrder(ctx context.Context, req types.CloseOrderRequest) (*types.CloseOrderResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paySig := s.client.GetPaySig(closeOrderURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := closeOrderURI + "?" + query.Encode()

	resp, err := s.client.Do(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("wechat api http status %d: %s", resp.StatusCode, string(raw))
	}

	var out types.CloseOrderResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// GetOrder 查询订单。
func (s *orderService) GetOrder(ctx context.Context, req types.GetOrderRequest) (*types.GetOrderResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paySig := s.client.GetPaySig(getOrderURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := getOrderURI + "?" + query.Encode()

	resp, err := s.client.Do(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("wechat api http status %d: %s", resp.StatusCode, string(raw))
	}

	var out types.GetOrderResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}
