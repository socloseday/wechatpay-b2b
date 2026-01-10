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

// BalanceService 处理余额查询与提现流程。
type BalanceService interface {
	// GetBalance 查询账户余额。
	GetBalance(ctx context.Context, req types.BalanceRequest) (*types.BalanceResponse, error)
	// Withdraw 发起提现。
	Withdraw(ctx context.Context, req types.WithdrawRequest) (*types.WithdrawResponse, error)
	// QueryWithdraw 查询提现状态。
	QueryWithdraw(ctx context.Context, req types.QueryWithdrawRequest) (*types.QueryWithdrawResponse, error)
}

type balanceService struct {
	client *client.Client
}

const (
	getMchBalanceURI = "/retail/B2b/getmchbalance"
	withdrawURI       = "/retail/B2b/withdraw"
	queryWithdrawURI  = "/retail/B2b/querywithdraw"
)

// NewBalanceService 创建余额服务。
func NewBalanceService(c *client.Client) BalanceService {
	return &balanceService{client: c}
}

// GetBalance 查询账户余额。
func (s *balanceService) GetBalance(ctx context.Context, req types.BalanceRequest) (*types.BalanceResponse, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
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

	paySig := s.client.GetPaySig(getMchBalanceURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := getMchBalanceURI + "?" + query.Encode()

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

	var out types.BalanceResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// Withdraw 发起提现。
func (s *balanceService) Withdraw(ctx context.Context, req types.WithdrawRequest) (*types.WithdrawResponse, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.WithdrawAmount <= 0 {
		return nil, errors.New("withdraw_amount is required")
	}
	if req.OutWithdrawNo == "" {
		return nil, errors.New("out_withdraw_no is required")
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

	paySig := s.client.GetPaySig(withdrawURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := withdrawURI + "?" + query.Encode()

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

	var out types.WithdrawResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}

// QueryWithdraw 查询提现状态。
func (s *balanceService) QueryWithdraw(ctx context.Context, req types.QueryWithdrawRequest) (*types.QueryWithdrawResponse, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("client is nil")
	}
	if req.Mchid == "" {
		return nil, errors.New("mchid is required")
	}
	if req.OutWithdrawNo == "" {
		return nil, errors.New("out_withdraw_no is required")
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

	paySig := s.client.GetPaySig(queryWithdrawURI, body)
	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	query.Set("pay_sig", paySig)
	uri := queryWithdrawURI + "?" + query.Encode()

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

	var out types.QueryWithdrawResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}
