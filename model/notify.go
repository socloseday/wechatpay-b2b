package model

// PaymentNotify 表示支付成功通知的 JSON 结构体。
type PaymentNotify struct {
	ToUserName         string        `json:"ToUserName"`           // 公众号/小程序ID
	FromUserName       string        `json:"FromUserName"`         // 用户OpenID
	CreateTime         int64         `json:"CreateTime"`           // 创建时间
	MsgType            string        `json:"MsgType"`              // 消息类型
	Event              string        `json:"Event"`                // 事件类型
	AppID              string        `json:"appid"`                // 小程序ID
	Mchid              string        `json:"mchid"`                // 微信商户号
	OutTradeNo         string        `json:"out_trade_no"`         // 商户订单号
	OrderID            string        `json:"order_id"`             // B2b支付订单号
	PayStatus          string        `json:"pay_status"`           // 支付状态
	PayTime            string        `json:"pay_time"`             // 支付完成时间
	Attach             string        `json:"attach"`               // 附加数据
	PayerOpenID        string        `json:"payer_openid"`         // 支付者OpenID
	Amount             PaymentAmount `json:"amount"`               // 订单金额信息
	WxPayTransactionID string        `json:"wxpay_transaction_id"` // 微信支付订单号
	Env                int           `json:"env"`                  // 订单环境
}

// PaymentAmount 表示通知中的金额信息。
type PaymentAmount struct {
	OrderAmount int64  `json:"order_amount"` // 订单金额
	PayerAmount int64  `json:"payer_amount"` // 支付者实付金额
	Currency    string `json:"currency"`     // 货币类型
}

// RefundNotify 表示退款通知的 JSON 结构体。
type RefundNotify struct {
	ToUserName         string `json:"ToUserName"`           // 公众号/小程序ID
	FromUserName       string `json:"FromUserName"`         // 用户OpenID
	CreateTime         int64  `json:"CreateTime"`           // 创建时间
	MsgType            string `json:"MsgType"`              // 消息类型
	Event              string `json:"Event"`                // 事件类型
	AppID              string `json:"appid"`                // 小程序ID
	Mchid              string `json:"mchid"`                // 微信商户号
	OutTradeNo         string `json:"out_trade_no"`         // 商户订单号
	OrderID            string `json:"order_id"`             // B2b支付订单号
	OutRefundNo        string `json:"out_refund_no"`        // 商户退款单号
	RefundID           string `json:"refund_id"`            // 微信退款订单号
	RefundStatus       string `json:"refund_status"`        // 退款状态
	RefundTime         string `json:"refund_time"`          // 退款完成时间
	PayerAmount        int64  `json:"payer_amount"`         // 退款金额
	WxPayTransactionID string `json:"wxpay_transaction_id"` // 微信支付订单号
	Env                int    `json:"env"`                  // 订单环境
}
