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

// ProfitService 处理分账相关调用。
type ProfitService interface {
	// ProfitSharing 请求分账。
	ProfitSharing(ctx context.Context, req types.ProfitSharingRequest) (*types.ProfitSharingResponse, error)
	// QueryProfitSharing 查询分账订单。
	QueryProfitSharing(ctx context.Context, req types.QueryProfitSharingRequest) (*types.QueryProfitSharingResponse, error)
	// ProfitSharingFinish 分账完结。
	ProfitSharingFinish(ctx context.Context, req types.ProfitSharingFinishRequest) (*types.ProfitSharingFinishResponse, error)
	// ProfitSharingReturn 分账回退。
	ProfitSharingReturn(ctx context.Context, req types.ProfitSharingReturnRequest) (*types.ProfitSharingReturnResponse, error)
	// QueryProfitSharingReturn 查询分账回退。
	QueryProfitSharingReturn(ctx context.Context, req types.QueryProfitSharingReturnRequest) (*types.QueryProfitSharingReturnResponse, error)
}

type profitService struct {
	client *client.Client
}

const (
	profitSharingURI         = "/retail/B2b/profitsharing"
	queryProfitSharingURI    = "/retail/B2b/queryprofitsharing"
	profitSharingFinishURI   = "/retail/B2b/profitsharingfinish"
	profitSharingReturnURI   = "/retail/B2b/profitsharingreturn"
	queryProfitSharingReturn = "/retail/B2b/queryprofitsharingreturn"
)

// NewProfitService 创建分账服务。
func NewProfitService(c *client.Client) ProfitService {
	return &profitService{client: c}
}

// ProfitSharing 请求分账。
func (s *profitService) ProfitSharing(ctx context.Context, req types.ProfitSharingRequest) (*types.ProfitSharingResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.TransactionID == "" {
		return nil, errors.New("transaction_id is required")
	}
	if req.OutOrderNo == "" {
		return nil, errors.New("out_order_no is required")
	}
	if len(req.Receivers) == 0 {
		return nil, errors.New("receivers is required")
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

	paySig := s.client.GetPaySig(profitSharingURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := profitSharingURI + "?" + query.Encode()

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

	var out types.ProfitSharingResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// QueryProfitSharing 查询分账订单。
func (s *profitService) QueryProfitSharing(ctx context.Context, req types.QueryProfitSharingRequest) (*types.QueryProfitSharingResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.TransactionID == "" {
		return nil, errors.New("transaction_id is required")
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

	paySig := s.client.GetPaySig(queryProfitSharingURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := queryProfitSharingURI + "?" + query.Encode()

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

	var out types.QueryProfitSharingResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// ProfitSharingFinish 分账完结。
func (s *profitService) ProfitSharingFinish(ctx context.Context, req types.ProfitSharingFinishRequest) (*types.ProfitSharingFinishResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.TransactionID == "" {
		return nil, errors.New("transaction_id is required")
	}
	if req.OutOrderNo == "" {
		return nil, errors.New("out_order_no is required")
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

	paySig := s.client.GetPaySig(profitSharingFinishURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := profitSharingFinishURI + "?" + query.Encode()

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

	var out types.ProfitSharingFinishResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// ProfitSharingReturn 分账回退。
func (s *profitService) ProfitSharingReturn(ctx context.Context, req types.ProfitSharingReturnRequest) (*types.ProfitSharingReturnResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.TransactionID == "" {
		return nil, errors.New("transaction_id is required")
	}
	if req.OutOrderNo == "" {
		return nil, errors.New("out_order_no is required")
	}
	if req.OutReturnNo == "" {
		return nil, errors.New("out_return_no is required")
	}
	if req.ReturnAmount <= 0 {
		return nil, errors.New("return_amount is required")
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

	paySig := s.client.GetPaySig(profitSharingReturnURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := profitSharingReturnURI + "?" + query.Encode()

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

	var out types.ProfitSharingReturnResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// QueryProfitSharingReturn 查询分账回退。
func (s *profitService) QueryProfitSharingReturn(ctx context.Context, req types.QueryProfitSharingReturnRequest) (*types.QueryProfitSharingReturnResponse, error) {
	if s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.OutReturnNo == "" && req.ReturnID == "" {
		return nil, errors.New("out_return_no or return_id is required")
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

	paySig := s.client.GetPaySig(queryProfitSharingReturn, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := queryProfitSharingReturn + "?" + query.Encode()

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

	var out types.QueryProfitSharingReturnResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}
