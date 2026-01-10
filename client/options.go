package client

import (
	"errors"
)

// Options Client 初始化参数。
type Options struct {
	// BaseURL 微信接口域名，默认 https://api.weixin.qq.com。
	BaseURL string
	// TokenProvider 接口调用凭证 access_token（可由业务侧定时刷新并更新到 Client 上）。
	TokenProvider string
	// AppKeyProvider 计算 pay_sig 所需的 appKey（建议按商户号维度维护不同 Client 实例）。
	AppKeyProvider string
}

// NewClient 创建一个可复用的微信 API Client。
// 注意：业务侧可定时刷新 access_token / appKey，并更新到返回的 Client 上。
func NewClient(opts Options) (*Client, error) {
	if opts.BaseURL == "" {
		opts.BaseURL = "https://api.weixin.qq.com"
	}
	if opts.TokenProvider == "" {
		return nil, errors.New("tokenProvider is empty")
	}
	return &Client{
		BaseURL:        opts.BaseURL,
		TokenProvider:  opts.TokenProvider,
		AppKeyProvider: opts.AppKeyProvider,
	}, nil
}
