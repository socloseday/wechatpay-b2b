package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
)

// Client 封装访问微信 API 的 HTTP 客户端。
type Client struct {
	baseURL     string
	accessToken string // access_token
	appKey      string // appKey（用于计算 pay_sig）
	httpClient  *http.Client
}

// GetAccessToken 获取 access_token。
func (c *Client) GetAccessToken() string {
	return c.accessToken
}

// GetAppKey 获取 appKey。
func (c *Client) GetAppKey() string {
	return c.appKey
}

// Do 向指定 uri（路径，不含完整域名）发起 HTTP 请求。
func (c *Client) Do(ctx context.Context, method, uri string, body []byte) (*http.Response, error) {
	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return httpClient.Do(req)
}

// GetPaySig 计算 pay_sig，算法为 HMAC-SHA256(appKey, uri+"&"+body)。
func (c *Client) GetPaySig(uri string, body []byte) string {
	return GetPaySig(uri, body, c.appKey)
}

// GetUserSignature 计算用户态签名，算法为 HMAC-SHA256(sessionKey, body)。
func (c *Client) GetUserSignature(body []byte, sessionKey string) string {
	return GetUserSignature(body, sessionKey)
}

// BuildURIWithAuth 构建带 access_token 参数的 URI。
func (c *Client) BuildURIWithAuth(uri string) string {
	query := url.Values{}
	query.Set("access_token", c.accessToken)
	return uri + "?" + query.Encode()
}

// BuildURIWithAuthAndSig 构建带 access_token 和 pay_sig 参数的 URI。
func (c *Client) BuildURIWithAuthAndSig(uri string, body []byte) string {
	paySig := c.GetPaySig(uri, body)
	query := url.Values{}
	query.Set("access_token", c.accessToken)
	query.Set("pay_sig", paySig)
	return uri + "?" + query.Encode()
}

// GetPaySig 计算 pay_sig，算法为 HMAC-SHA256(appKey, uri+"&"+body)。
func GetPaySig(uri string, body []byte, appKey string) string {
	msg := uri + "&" + string(body)
	return hmacHex(appKey, msg)
}

// GetUserSignature 计算用户态签名，算法为 HMAC-SHA256(sessionKey, body)。
func GetUserSignature(body []byte, sessionKey string) string {
	return hmacHex(sessionKey, string(body))
}

func hmacHex(key, msg string) string {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
