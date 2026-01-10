# wechatpay-b2b

微信支付 B2B（门店助手）Golang SDK。

## 安装

```bash
go get github.com/enjoy322/wechatpay-b2b
```

## 项目定位

- 面向小程序 B2B 门店助手的"支付/退款/提现/合单支付"等能力，提供后端 SDK 封装（签名、请求结构体、调用入口等）。
- 小程序侧支付入口为 `wx.requestCommonPayment`，本 SDK 负责在服务端准备 `signData`、`paySig` 等参数并对接相关服务端接口。

## 官方文档

- 文档：<https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/industry/B2b_store_assistant.html>

## 支持的接口

### 订单服务 (OrderService)

| 方法 | 功能 | URI |
|-----|-----|-----|
| `CloseOrder` | 关闭订单 | `/retail/B2b/closeb2border` |
| `GetOrder` | 查询订单 | `/retail/B2b/getorder` |

### 退款服务 (RefundService)

| 方法 | 功能 | URI |
|-----|-----|-----|
| `CreateRefund` | 发起退款 | `/retail/B2b/createrefund` |
| `GetRefund` | 查询退款 | `/retail/B2b/getrefund` |

### 余额服务 (BalanceService)

| 方法 | 功能 | URI |
|-----|-----|-----|
| `GetBalance` | 查询账户余额 | `/retail/B2b/getmchbalance` |
| `Withdraw` | 发起提现 | `/retail/B2b/withdraw` |
| `QueryWithdraw` | 查询提现状态 | `/retail/B2b/querywithdraw` |

### 分账服务 (ProfitService)

| 方法 | 功能 | URI |
|-----|-----|-----|
| `ProfitSharing` | 请求分账 | `/retail/B2b/profitsharing` |
| `QueryProfitSharing` | 查询分账订单 | `/retail/B2b/queryprofitsharing` |
| `ProfitSharingFinish` | 分账完结 | `/retail/B2b/profitsharingfinish` |
| `ProfitSharingReturn` | 分账回退 | `/retail/B2b/profitsharingreturn` |
| `QueryProfitSharingReturn` | 查询分账回退 | `/retail/B2b/queryprofitsharingreturn` |

### 门店服务 (RetailService)

| 方法 | 功能 | URI |
|-----|-----|-----|
| `BatchCreateRetail` | 预录入门店信息 | `/wxa/business/batchcreateretail` |

### 通知解析

| 函数 | 功能 |
|-----|-----|
| `model.ParsePaymentNotify` | 解析支付通知 |
| `model.ParseRefundNotify` | 解析退款通知 |

## 使用说明

### Client 配置（access_token / appKey）

本 SDK 以 `client.Client` 作为共享调用上下文：

- `access_token`：保存在 `client.Client.AccessToken`，业务侧定时刷新并更新即可。
- `appKey`：保存在 `client.Client.AppKey`，用于服务端接口计算 `pay_sig`（建议按商户号维度维护不同的 `client.Client` 实例）。

因此调用 `service.BalanceService.GetBalance` 等方法时，不需要每次显式传入 `access_token` / `appKey`，只需确保 `client.Client` 内的相关字段是最新的。

建议直接替换整个 `client.Client` 实例。

### 示例

```go
package main

import (
    "context"
    "fmt"

    "github.com/enjoy322/wechatpay-b2b/client"
    "github.com/enjoy322/wechatpay-b2b/service"
    "github.com/enjoy322/wechatpay-b2b/types"
)

func main() {
    // 初始化客户端
    c := client.NewClient(client.Options{
        BaseURL:       "https://api.weixin.qq.com",
        TokenProvider: "your_access_token",
        AppKeyProvider: "your_app_key",
    })

    // 创建服务
    orderSvc := service.NewOrderService(c)

    // 查询订单
    resp, err := orderSvc.GetOrder(context.Background(), types.GetOrderRequest{
        Mchid:      "1230000109",
        OutTradeNo: "your_out_trade_no",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Order: %+v\n", resp)
}
```
