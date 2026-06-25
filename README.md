# go-aster

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-aster/v3.svg)](https://pkg.go.dev/github.com/UnipayFI/go-aster/v3)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go SDK for the [Aster DEX](https://www.asterdex.com) **V3** API (Spot + Futures, REST + WebSocket).

> Aster V3 replaces the legacy HMAC-SHA256 auth with an **API-wallet model** based on EIP-712 typed data and ECDSA signatures. The legacy V1 API stopped issuing new keys on 2026-03-25; existing keys still work, but ongoing development happens here on V3. If you need V1, switch to the `V1(Legacy)` branch.

| API                                              | Aligned to                                                                                          |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------- |
| Spot + Futures V3 REST + WebSocket (public / private) | [2026-6-17](https://asterdex.github.io/aster-api-website/changelog/upcoming-changes/) |

---

## Features

- ✅ **Spot REST**: 24 endpoints — market data, orders, account, transfers, withdrawals
- ✅ **Futures REST**: ~55 endpoints — market, trading (incl. chase & strategy orders), position, account, STP, MMP, sub-accounts, asset migration, agent registration, announcements
- ✅ **WebSocket Streaming**: full Spot & Futures market streams plus user data streams
- ✅ **Sub-account Flows**: bind / create / update / transfer with master + child signature inputs
- ✅ **Flexible Signer**: pluggable `SignFn` for HSM, TEE, or remote signing
- ✅ **TEE Injection**: `WithTEEAuth` lets the signer's private key stay inside the enclave — the SDK only holds the user/signer addresses and delegates signing via `WithSignRequestFn`
- ✅ **Proxy Support**: `WithProxy` / `WithWebSocketProxy` route REST and WebSocket traffic through HTTP, HTTPS, or SOCKS5 (`socks5` / `socks5h`) proxies, with optional `user:pass@` auth in the URL
- ✅ **Mainnet & Testnet**: switch with `WithNetwork(common.Mainnet|common.Testnet)` — base URL and chainId both follow

---

## 📦 Installation

```bash
go get -u github.com/UnipayFI/go-aster/v3
```

> Module path is `github.com/UnipayFI/go-aster/v3` (Go major-version convention). Requires **Go 1.21+**.

---

## 🚀 Quick Start

### 1. Public market data (no auth)

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

### 2. Authenticated client

```go
import (
    aster "github.com/UnipayFI/go-aster/v3"
    "github.com/UnipayFI/go-aster/v3/client"
)

c := aster.NewSpotClient(
    client.WithAuth(
        "0x...your master wallet address...",
        "0x...your API wallet PRIVATE KEY (hex)...",
    ),
)
```

`WithAuth` only needs the **API wallet private key** — the signer address is derived automatically. The first argument is the **master wallet address**; its private key never touches the client.

### 3. Top-level constructors

```go
aster.NewSpotClient(...)              // /api/v3/* REST
aster.NewSpotWebSocketClient(...)     // wss://sstream.asterdex.com
aster.NewFuturesClient(...)           // /fapi/v3/* REST
aster.NewFuturesWebSocketClient(...)  // wss://fstream.asterdex.com
```

### 4. Configuration options

```go
client.WithAuth(userAddress, signerPrivateKeyHex string)
client.WithTEEAuth(userAddress, signerAddress string) // TEE/HSM mode: no local private key
client.WithNetwork(common.Network)      // common.Mainnet (default) / common.Testnet — sets base URL + chainId
client.WithBaseURL(string)              // override the network's REST base URL (chainId still follows WithNetwork)
client.WithLogger(log.Logger)           // slog-compatible interface
client.WithSignRequestFn(client.SignFn) // custom signer (HSM / TEE / remote)
client.WithRecvWindow(int64)
client.WithTimeOffset(int64)            // ms; aligns nonce with server clock
```

For WebSocket clients, use `client.WithWebSocketNetwork(common.Network)` and `client.WithWebSocketBaseURL(string)` (the base URL override likewise takes priority over the network default).

---

## 📖 Usage Examples

### Place a Spot limit order

```go
import (
    "github.com/UnipayFI/go-aster/v3/spot"
    "github.com/shopspring/decimal"
)

order, err := c.NewPlaceOrderService("BTCUSDT", spot.SideBuy, spot.OrderTypeLimit).
    SetTimeInForce(spot.TimeInForceGTC).
    SetQuantity(decimal.NewFromFloat(0.01)).
    SetPrice(decimal.NewFromInt(60000)).
    Do(context.Background())
```

### Subscribe to a market WebSocket stream

```go
ws := aster.NewSpotWebSocketClient()

done, stop, err := ws.NewSubscribeAggTradeService("BTCUSDT").
    Do(context.Background(), func(ev *spot.WsAggTradeEvent, err error) {
        if err != nil {
            return
        }
        fmt.Println(ev.Symbol, ev.Price, ev.Quantity)
    })
_, _, _ = done, stop, err
```

### Subscribe to the user data stream

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

### TEE injection (private key never leaves the enclave)

For deployments where the API-wallet private key must stay inside a trusted execution environment, use `WithTEEAuth` to register the user/signer addresses and delegate the actual signing to your TEE binary via `WithSignRequestFn`. The SDK never sees the private key.

```go
import (
    "os/exec"
    "strings"

    aster "github.com/UnipayFI/go-aster/v3"
    "github.com/UnipayFI/go-aster/v3/client"
)

c := aster.NewFuturesClient(
    client.WithTEEAuth(
        "0x...master wallet address...",
        "0x...API wallet (signer) address...",
    ),
    client.WithSignRequestFn(func(_ /*privateKeyHex*/, msg string, _ /*chainID*/ int64) (string, error) {
        // Shell out to the TEE binary. API_KEY is read from the env by the TEE
        // process; the private key is sealed inside the enclave.
        out, err := exec.Command("./tee",
            "--encrypt-type=EIP712-ECDSA",
            "--input="+msg,
        ).Output()
        if err != nil {
            return "", err
        }
        return strings.TrimSpace(string(out)), nil
    }),
)
```

The companion TEE binary lives at [`UnipayFI/tee-encrypt`](https://github.com/UnipayFI/tee-encrypt). It reads the sealed signer key via the `API_KEY` env var, computes the EIP-712 digest internally (chainId `1666`, the protocol-fixed Aster domain) and returns a 130-char hex signature with `v ∈ {27, 28}` — byte-compatible with the SDK's local signer.

### Sub-account flows (master / child signatures)

Sub-account endpoints (`Bind`, `Create`, `Update`, `Transfer`) require signatures from the **master wallet** — and sometimes the **child wallet** — rather than the signer/agent. Because master keys typically live in cold storage, those services accept already-computed signatures as inputs:

```go
import "github.com/UnipayFI/go-aster/v3/request"

// Caller computes signatures via request.SignEIP712V3 with the appropriate
// message body (see godoc on each service for the exact format).
sig, _ := request.SignEIP712V3(masterPrivKeyHex, msgBody, 1666)

c.NewBindSubAccountService(childAddr, name, user, nonce, childSig, sig).
    Do(ctx)
```

`request.SignEIP712V3(privateKeyHex, msg, chainID)` and `request.EIP712Digest(msg, chainID)` are exposed for callers that need to interact with these flows or implement their own signers.

### Signing internals

V3 signing uses a fixed EIP-712 typed-data envelope:

- `domain.name = "AsterSignTransaction"`, `version = "1"`, `verifyingContract = 0x000...0`
- `chainId = 1666` (mainnet) / `714` (testnet) — derived from `WithNetwork`
- `primaryType = "Message"`, single field `msg` of type `string`
- `msg` value = the URL-encoded query string of the request (already including `nonce` and `signer`)
- Output signature has `v` adjusted to `27/28` to match `eth_account.sign_message`

See `request/sign.go` for the implementation and `request/sign_test.go` for the regression tests (digest stability, ecrecover round-trip, deterministic signatures, chainId sensitivity).

### Endpoint coverage

<details>
<summary><b>Spot REST (24)</b></summary>

`Ping`, `GetServerTime`, `Noop`, `GetExchangeInfo`, `GetDepth`, `GetRecentTrades`, `GetHistoricalTrades`, `GetAggTrades`, `GetKlines`, `Get24hTicker`, `GetTickerPrice`, `GetBookTicker`, `GetCommissionRate`, `PlaceOrder`, `CancelOrder`, `GetOrder`, `GetOpenOrder`, `GetOpenOrders`, `CancelAllOpenOrders`, `GetAllOrders`, `GetTransactionHistory`, `PerpSpotTransfer`, `GetWithdrawFee`, `Withdraw`, `GetAccount`, `GetUserTrades`.

</details>

<details>
<summary><b>Spot WebSocket</b></summary>

`Subscribe{AggTrade,Trade,Kline,MiniTicker,AllMiniTickers,Ticker,AllTickers,BookTicker,AllBookTickers,PartialDepth,DiffDepth,TradePro,UserDataStream}`. ListenKey REST: `Create/Renew/DeleteListenKey`.

</details>

<details>
<summary><b>Futures REST (~45)</b></summary>

- **General**: `Ping`, `Time`, `Noop`
- **Market data**: `ExchangeInfo`, `Depth`, `Trades`, `HistoricalTrades`, `AggTrades`, `Klines`, `IndexPriceKlines`, `MarkPriceKlines`, `PremiumIndex`, `FundingRate`, `FundingInfo`, `24hTicker`, `TickerPrice`, `BookTicker`, `IndexReferences`
- **Trading**: `PlaceOrder`, `ModifyOrder`, `BatchOrders`, `FuturesSpotTransfer`, `GetOrder`, `CancelOrder`, `CancelAllOpenOrders`, `CancelMultipleOrders`, `CountdownCancelAll`, `GetOpenOrder`, `GetOpenOrders`, `GetAllOrders`
- **Position config**: `ChangePositionMode`, `GetPositionMode`, `ChangeMultiAssetsMode`, `GetMultiAssetsMode`, `ChangeLeverage`, `ChangeMarginType`, `ModifyIsolatedPositionMargin`, `GetPositionMarginHistory`, `PositionRisk`, `ADLQuantile`, `ForceOrders`
- **Account**: `Balance`, `Account`, `UserTrades`, `IncomeHistory`, `LeverageBracket`, `CommissionRate`
- **MMP**: `UpdateMMP`, `GetMMP`, `DeleteMMP`, `ResetMMP`
- **Sub-accounts**: `Bind`, `Create`, `GetList`, `Update`, `Transfer`

</details>

<details>
<summary><b>Futures WebSocket</b></summary>

- **Market**: `AggTrade`, `MarkPrice`, `AllMarkPrices`, `Kline`, `MiniTicker`, `AllMiniTickers`, `Ticker`, `AllTickers`, `BookTicker`, `AllBookTickers`, `ForceOrder`, `AllForceOrders`, `PartialDepth`, `DiffDepth`
- **User data**: `UserDataStream`
- **ListenKey REST**: `Create/Renew/DeleteListenKey`

</details>

---

## 📄 License

Released under the [MIT License](LICENSE).
