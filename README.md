# go-aster

A comprehensive Go SDK for AsterDEX API, supporting both Spot and Futures trading with REST API and WebSocket streaming capabilities.

## Features

- ✅ **Complete API Coverage**: Full support for Spot and Futures trading
- ✅ **WebSocket Streaming**: Real-time market data and user data streams
- ✅ **Flexible Authentication**: Built-in signature support with custom signature function capability (including TEE support)
- ✅ **Type Safety**: Strongly typed responses and parameters
- ☑️ **V3 API Support**: Spot&Futures V3 Support
- ☑️ **Deposit Withdrawal API Support**: Deposit&Withdrawal API Support

## Installation

```bash
go get -u github.com/UnipayFI/go-aster
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
)

func main() {
    // Initialize clients with API credentials
    apiKey := "your_api_key"
    secretKey := "your_secret_key"
    
    // Create spot client
    spotClient := aster.NewSpotClient(
        client.WithAuth(apiKey, secretKey),
    )
    
    // Create futures client
    futuresClient := aster.NewFuturesClient(
        client.WithAuth(apiKey, secretKey),
    )
}
```

## Spot Trading Examples

### Market Data Queries

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
)

func main() {
    client := aster.NewSpotClient()
    ctx := context.Background()

    // Get market depth
    depth, err := client.NewDepthService("BTCUSDT").
        SetLimit(10).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Best Bid: %s, Best Ask: %s\n", 
        depth.Bids[0].Price(), depth.Asks[0].Price())

    // Get recent trades
    trades, err := client.NewGetTradesService("BTCUSDT").
        SetLimit(5).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest trade price: %s\n", trades[0].Price)

    // Get exchange info
    exchangeInfo, err := client.NewExchangeInfoService().
        GetExchangeInfo(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Server time: %s\n", exchangeInfo.ServerTime)
}
```

### Account Information

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
)

func main() {
    client := aster.NewSpotClient(
        client.WithAuth("your_api_key", "your_secret_key"),
    )
    ctx := context.Background()

    // Sync server time (recommended before trading)
    if err := client.SyncServerTime(ctx); err != nil {
        log.Fatal(err)
    }

    // Get account information
    account, err := client.NewGetAccountService().Do(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Can trade: %t\n", account.CanTrade)
    for _, balance := range account.Balances {
        if balance.Free.IsPositive() || balance.Locked.IsPositive() {
            fmt.Printf("Asset: %s, Free: %s, Locked: %s\n", 
                balance.Asset, balance.Free, balance.Locked)
        }
    }
}
```

### Spot Order Management

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
    "github.com/UnipayFI/go-aster/spot"
)

func main() {
    client := aster.NewSpotClient(
        client.WithAuth("your_api_key", "your_secret_key"),
    )
    ctx := context.Background()

    // Sync server time
    if err := client.SyncServerTime(ctx); err != nil {
        log.Fatal(err)
    }

    // Place a limit buy order
    buyOrder, err := client.NewCreateOrderService("BTCUSDT", spot.OrderSideBuy, spot.OrderTypeLimit).
        SetQuantity(0.001).
        SetPrice(30000.0).
        SetTimeInForce(spot.TimeInForceTypeGTC).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Buy order placed: ID %d, Status: %s\n", 
        buyOrder.OrderId, buyOrder.Status)

    // Get order status
    orderStatus, err := client.NewGetOrderService("BTCUSDT").
        SetOrderId(buyOrder.OrderId).
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order status: %s, Executed: %s\n", 
        orderStatus.Status, orderStatus.ExecutedQty)

    // Get all open orders
    openOrders, err := client.NewGetOpenOrdersService().
        SetSymbol("BTCUSDT").
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Open orders count: %d\n", len(openOrders))

    // Cancel the order (if still open)
    if orderStatus.Status == spot.OrderStatusNew {
        canceledOrder, err := client.NewCancelOrderService("BTCUSDT").
            SetOrderId(buyOrder.OrderId).
            Do(ctx)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Order canceled: %s\n", canceledOrder.Status)
    }
}
```

## Futures Trading Examples

### Futures Market Data and Orders

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
    "github.com/UnipayFI/go-aster/futures"
)

func main() {
    client := aster.NewFuturesClient(
        client.WithAuth("your_api_key", "your_secret_key"),
    )
    ctx := context.Background()

    // Sync server time
    if err := client.SyncServerTime(ctx); err != nil {
        log.Fatal(err)
    }

    // Get futures account balance
    account, err := client.NewGetAccountService().Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total wallet balance: %s USDT\n", account.TotalWalletBalance)

    // Place a futures order
    order, err := client.NewCreateOrderService("BTCUSDT", futures.OrderSideBuy, futures.OrderTypeLimit).
        SetQuantity(0.001).
        SetPrice(30000.0).
        SetTimeInForce(futures.TimeInForceTypeGTC).
        SetPositionSide(futures.PositionSideLong). // For hedge mode
        Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Futures order placed: ID %d\n", order.OrderId)

    // Get position information
    positions, err := client.NewGetPositionRiskService().Do(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, pos := range positions {
        if pos.PositionAmt.IsPositive() {
            fmt.Printf("Position: %s, Size: %s, PnL: %s\n", 
                pos.Symbol, pos.PositionAmt, pos.UnRealizedProfit)
        }
    }
}
```

## WebSocket Examples

### Public Market Data Streaming

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
    "github.com/UnipayFI/go-aster/spot"
)

func main() {
    wsClient := aster.NewSpotWebSocketClient()
    ctx := context.Background()

    // Subscribe to depth updates
    done, stop, err := wsClient.NewSubscribeDepthService("BTCUSDT", "5").
        Do(ctx, func(event *spot.WsDepthResponse, err error) {
            if err != nil {
                log.Printf("Error: %v", err)
                return
            }
            fmt.Printf("Depth update - Symbol: %s, Bids: %d, Asks: %d\n", 
                event.Symbol, len(event.Bids), len(event.Asks))
        })
    if err != nil {
        log.Fatal(err)
    }
    select {}
}
```

### User Data Streaming

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
    "github.com/UnipayFI/go-aster/spot"
)

// Implement the UserDataHandler interface
type MyUserDataHandler struct{}

func (h *MyUserDataHandler) OnAccountUpdate(event *spot.WsAccountUpdateEvent) {
    fmt.Printf("Account update - Event time: %s\n", event.EventTime)
}

func (h *MyUserDataHandler) OnOrderUpdate(event *spot.WsOrderUpdateEvent) {
    fmt.Printf("Order update - Symbol: %s, OrderId: %d, Status: %s, Side: %s\n",
        event.Symbol, event.OrderId, event.Status, event.Side)
}

func (h *MyUserDataHandler) OnListenKeyExpired(event *spot.WsListenKeyExpiredEvent) {
    fmt.Printf("Listen key expired at: %s\n", event.EventTime)
}

func (h *MyUserDataHandler) OnError(err error) {
    log.Printf("WebSocket error: %v", err)
}

func main() {
    // Create REST client to get listen key
    restClient := aster.NewSpotClient(
        client.WithAuth("your_api_key", "your_secret_key"),
    )
    
    // Create WebSocket client
    wsClient := aster.NewSpotWebSocketClient()
    ctx := context.Background()

    // Get listen key for user data stream
    listenKey, err := restClient.NewCreateListenKeyService().Do(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Subscribe to user data stream
    handler := &MyUserDataHandler{}
    done, stop, err := wsClient.NewSubscribeUserDataStreamService(listenKey).
        Do(ctx, handler)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Listening for user data updates... Press Ctrl+C to stop")

    // Wait for interrupt signal
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    <-interrupt
    log.Println("Shutting down...")
    
    done <- struct{}{}
    <-stop
    
    log.Println("Connection closed")
}
```

## Advanced Configuration

### Custom Signature Function (TEE Support)

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"

    "github.com/UnipayFI/go-aster"
    "github.com/UnipayFI/go-aster/client"
)

// Example custom signature function for TEE integration
func customTEESignature(apiKey, secretKey string, payload string) (string, error) {
    // This is where you would integrate with your TEE solution
    // For example, call Intel SGX enclave or remote HSM
    
    // Fallback to standard HMAC-SHA256 for this example
    mac := hmac.New(sha256.New, []byte(secretKey))
    mac.Write([]byte(payload))
    return fmt.Sprintf("%x", mac.Sum(nil)), nil
}

func main() {
    client := aster.NewSpotClient(
        client.WithAuth("your_api_key", ""), // Protect API_SECRET
        client.WithSignRequestFn(customTEESignature), // Custom signature function
    )
    
    // Client will now use your custom signature function for all signed requests
}
```

## Error Handling

The SDK provides structured error handling:

```go
order, err := c.NewCreateOrderService("BTCUSDT", spot.OrderSideBuy, spot.OrderTypeLimit).
    SetQuantity(0.001).
    SetPrice(30000.0).
    Do(ctx)
	if err != nil {
		if err, ok := err.(*client.APIError); ok {
			logrus.Errorf("error code: %d, error message: %s", err.Code, err.Message)
			return
		}
		logrus.Errorf("create order error: %v", err)
		return
	}
```

## Rate Limiting

The SDK automatically handles rate limiting information returned by the API:

```go
// Check rate limit usage
usedWeight := client.GetUsedWeight()
fmt.Printf("Used weight: %d/%d\n", usedWeight.Used, usedWeight.Used1M)

orderCount := client.GetOrderCount()
fmt.Printf("Orders: %d (10s), %d (1d)\n", orderCount.Count10s, orderCount.Count1d)
```

## Licence
UnipayFI/go-aster is released under the [MIT License](https://opensource.org/licenses/MIT).
