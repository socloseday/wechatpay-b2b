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

// RetailService 处理门店助手（wxa/business）相关接口调用。
type RetailService interface {
	// BatchCreateRetail 预录入门店信息。
	BatchCreateRetail(ctx context.Context, req types.BatchCreateRetailRequest) (*types.BatchCreateRetailResponse, error)
}

type retailService struct {
	client *client.Client
}

const batchCreateRetailURI = "/wxa/business/batchcreateretail"

// NewRetailService 创建门店助手服务。
func NewRetailService(c *client.Client) RetailService {
	return &retailService{client: c}
}

// BatchCreateRetail 预录入门店信息。
func (s *retailService) BatchCreateRetail(ctx context.Context, req types.BatchCreateRetailRequest) (*types.BatchCreateRetailResponse, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("client is nil")
	}
	if len(req.RetailInfoList) == 0 {
		return nil, errors.New("retail_info_list is required")
	}
	if s.client.GetAccessToken() == "" {
		return nil, errors.New("accessToken is empty")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("access_token", s.client.GetAccessToken())
	uri := batchCreateRetailURI + "?" + query.Encode()

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

	var out types.BatchCreateRetailResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return &out, fmt.Errorf("wechat api error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return &out, nil
}
