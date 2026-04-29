# go-aster

Go SDK for the [Aster DEX](https://www.asterdex.com) **V3** API (Spot + Futures, REST + WebSocket).

V3 uses **EIP-712 + ECDSA** signing (the API wallet model) instead of the legacy V1 HMAC-SHA256 scheme. The legacy V1 API is no longer accepting new keys (since 2026-03-25); existing V1 keys still work, but the V1 SDK is no longer maintained on `main` — see the `V1(Legacy)` branch if you need it.

## Install

```bash
go get -u github.com/UnipayFI/go-aster/v3
```

Module path is `github.com/UnipayFI/go-aster/v3` (Go major-version convention).

## Quick start

### REST — public market data (no auth)

```go
package main

import (
    "context"
    "fmt"

    aster "github.com/UnipayFI/go-aster/v3"
)

func main() {
    c := aster.NewSpotClient()
    info, err := c.NewGetExchangeInfoService().Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Println("symbols:", len(info.Symbols))
}
```

### REST — signed endpoints (place an order)

```go
import (
    aster "github.com/UnipayFI/go-aster/v3"
    "github.com/UnipayFI/go-aster/v3/client"
    "github.com/UnipayFI/go-aster/v3/spot"
    "github.com/shopspring/decimal"
)

c := aster.NewSpotClient(
    client.WithAuth(
        "0x...your main wallet address...",
        "0x...your API wallet PRIVATE KEY (hex)...",
    ),
)

order, err := c.NewPlaceOrderService("BTCUSDT", spot.SideBuy, spot.OrderTypeLimit).
    SetTimeInForce(spot.TimeInForceGTC).
    SetQuantity(decimal.NewFromFloat(0.01)).
    SetPrice(decimal.NewFromInt(60000)).
    Do(context.Background())
```

`WithAuth` only needs the API wallet's **private key** — the signer address is derived from it. The first argument (`user`) is the master wallet address; it does not need a private key on the client side.

For testnet, set `client.WithChainID(714)` and a testnet base URL via `client.WithBaseURL(...)`. Mainnet (chainId 1666) is the default.

### WebSocket — market stream

```go
ws := aster.NewSpotWebSocketClient()

done, stop, err := ws.NewSubscribeAggTradeService("BTCUSDT").
    Do(context.Background(), func(ev *spot.WsAggTradeEvent, err error) {
        if err != nil {
            return
        }
        fmt.Println(ev.Symbol, ev.Price, ev.Quantity)
    })
_ = done
_ = stop
_ = err
```

### WebSocket — user data stream

```go
spotREST := aster.NewSpotClient(client.WithAuth(user, signerPrivKeyHex))
key, _ := spotREST.NewCreateListenKeyService().Do(context.Background())

wsClient := aster.NewSpotWebSocketClient()
wsClient.NewSubscribeUserDataStreamService(key.ListenKey).
    Do(context.Background(), func(ev *spot.WsUserDataEvent, err error) {
        switch ev.EventType {
        case "executionReport":
            fmt.Println("order update:", ev.ExecutionReport.OrderStatus)
        case "outboundAccountPosition":
            fmt.Println("balances:", ev.AccountUpdate.Balances)
        }
    })
```

## Top-level constructors

```go
aster.NewSpotClient(...)               // /api/v3/* REST
aster.NewSpotWebSocketClient(...)      // wss://sstream.asterdex.com
aster.NewFuturesClient(...)            // /fapi/v3/* REST
aster.NewFuturesWebSocketClient(...)   // wss://fstream.asterdex.com
```

## Configuration

```go
client.WithAuth(userAddress, signerPrivateKeyHex string)
client.WithBaseURL(url string)         // override default endpoint
client.WithChainID(int64)              // 1666 mainnet (default), 714 testnet
client.WithLogger(log.Logger)          // slog-compatible interface
client.WithSignRequestFn(client.SignFn) // custom signing (HSM / TEE / remote)
client.WithRecvWindow(int64)
client.WithTimeOffset(int64)           // ms; aligns nonce with server clock
```

## Coverage

**Spot REST** (24): `NewPingService`, `NewGetServerTimeService`, `NewNoopService`, `NewGetExchangeInfoService`, `NewGetDepthService`, `NewGetRecentTradesService`, `NewGetHistoricalTradesService`, `NewGetAggTradesService`, `NewGetKlinesService`, `NewGet24hTickerService`, `NewGetTickerPriceService`, `NewGetBookTickerService`, `NewGetCommissionRateService`, `NewPlaceOrderService`, `NewCancelOrderService`, `NewGetOrderService`, `NewGetOpenOrderService`, `NewGetOpenOrdersService`, `NewCancelAllOpenOrdersService`, `NewGetAllOrdersService`, `NewGetTransactionHistoryService`, `NewPerpSpotTransferService`, `NewGetWithdrawFeeService`, `NewWithdrawService`, `NewGetAccountService`, `NewGetUserTradesService`.

**Spot WebSocket**: `NewSubscribe{AggTrade,Trade,Kline,MiniTicker,AllMiniTickers,Ticker,AllTickers,BookTicker,AllBookTickers,PartialDepth,DiffDepth,TradePro,UserDataStream}Service`. ListenKey REST: `NewCreate/Renew/DeleteListenKeyService`.

**Futures REST** (~45): general (`Ping`, `Time`, `Noop`), market data (`ExchangeInfo`, `Depth`, `Trades`, `HistoricalTrades`, `AggTrades`, `Klines`, `IndexPriceKlines`, `MarkPriceKlines`, `PremiumIndex`, `FundingRate`, `FundingInfo`, `24hTicker`, `TickerPrice`, `BookTicker`, `IndexReferences`), trading (`PlaceOrder`, `ModifyOrder`, `BatchOrders`, `FuturesSpotTransfer`, `GetOrder`, `CancelOrder`, `CancelAllOpenOrders`, `CancelMultipleOrders`, `CountdownCancelAll`, `GetOpenOrder`, `GetOpenOrders`, `GetAllOrders`), position config (`ChangePositionMode`, `GetPositionMode`, `ChangeMultiAssetsMode`, `GetMultiAssetsMode`, `ChangeLeverage`, `ChangeMarginType`, `ModifyIsolatedPositionMargin`, `GetPositionMarginHistory`, `PositionRisk`, `ADLQuantile`, `ForceOrders`), account (`Balance`, `Account`, `UserTrades`, `IncomeHistory`, `LeverageBracket`, `CommissionRate`), MMP (`UpdateMMP`, `GetMMP`, `DeleteMMP`, `ResetMMP`), sub-accounts (`Bind`, `Create`, `GetList`, `Update`, `Transfer`).

**Futures WebSocket**: market (`AggTrade`, `MarkPrice`, `AllMarkPrices`, `Kline`, `MiniTicker`, `AllMiniTickers`, `Ticker`, `AllTickers`, `BookTicker`, `AllBookTickers`, `ForceOrder`, `AllForceOrders`, `PartialDepth`, `DiffDepth`), user data (`UserDataStream`). ListenKey REST: `NewCreate/Renew/DeleteListenKeyService`.

## Sub-account flows

Sub-account endpoints (`Bind`, `Create`, `Update`, `Transfer`) require signatures from the **master wallet** (and sometimes the **child wallet**) private keys, not the signer/agent. Because master keys typically live in cold storage, those services accept already-computed signatures as inputs:

```go
// Caller computes signatures via request.SignEIP712V3 with the appropriate
// message body (see godoc on each service for the exact format).
sig, _ := request.SignEIP712V3(masterPrivKeyHex, msgBody, 1666)

c.NewBindSubAccountService(childAddr, name, user, nonce, childSig, sig).
    Do(ctx)
```

`request.SignEIP712V3(privateKeyHex, msg, chainID)` and `request.EIP712Digest(msg, chainID)` are exposed for callers that need to interact with these flows or implement their own signers.

## Signing internals

V3 signing uses fixed EIP-712 typed data:

- `domain.name = "AsterSignTransaction"`, `version = "1"`, `verifyingContract = 0x000...0`, `chainId = 1666` (mainnet) / `714` (testnet)
- `primaryType = "Message"`, single field `msg` of type `string`
- `msg` value = the URL-encoded query string of the request (already including `nonce` and `signer`)
- Output signature has `v` adjusted to `27/28` to match `eth_account.sign_message`

See `request/sign.go` for the implementation and `request/sign_test.go` for the regression tests (digest stability, ecrecover round-trip, deterministic signatures, chainId sensitivity).

## License

MIT — see [LICENSE](LICENSE).
