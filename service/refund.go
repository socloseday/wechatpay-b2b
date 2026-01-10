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

// RefundService 处理退款申请与退款查询。
type RefundService interface {
	// CreateRefund 发起退款。
	CreateRefund(ctx context.Context, req types.RefundRequest) (*types.RefundResponse, error)
	// GetRefund 查询退款。
	GetRefund(ctx context.Context, req types.GetRefundRequest) (*types.GetRefundResponse, error)
}

type refundService struct {
	client *client.Client
}

const (
	createRefundURI = "/retail/B2b/createrefund"
	getRefundURI    = "/retail/B2b/getrefund"
)

// NewRefundService 创建退款服务。
func NewRefundService(c *client.Client) RefundService {
	return &refundService{client: c}
}

// CreateRefund 发起退款。
func (s *refundService) CreateRefund(ctx context.Context, req types.RefundRequest) (*types.RefundResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.OutRefundNo == "" {
		return nil, errors.New("out_refund_no is required")
	}
	if req.RefundAmount <= 0 {
		return nil, errors.New("refund_amount is required")
	}
	if s.client.GetAccessToken() == "" {
		return nil, errors.New("accessToken is empty")
	}
	if s.client.GetAppKey() == "" {
		return nil, errors.New("appKey is empty")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paySig := s.client.GetPaySig(createRefundURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := createRefundURI + "?" + query.Encode()

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

	var out types.RefundResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// GetRefund 查询退款。
func (s *refundService) GetRefund(ctx context.Context, req types.GetRefundRequest) (*types.GetRefundResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.OutRefundNo == "" && req.RefundID == "" {
		return nil, errors.New("out_refund_no or refund_id is required")
	}
	if s.client.GetAccessToken() == "" {
		return nil, errors.New("accessToken is empty")
	}
	if s.client.GetAppKey() == "" {
		return nil, errors.New("appKey is empty")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	paySig := s.client.GetPaySig(getRefundURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := getRefundURI + "?" + query.Encode()

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

	var out types.GetRefundResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}
