package aster

import (
	"github.com/UnipayFI/go-aster/client"
	"github.com/UnipayFI/go-aster/futures"
	"github.com/UnipayFI/go-aster/spot"
)

func NewSpotClient(options ...client.Options) *spot.SpotClient {
	return spot.NewSpotClient(options...)
}

func NewSpotWebSocketClient(options ...client.Options) *spot.SpotWebSocketClient {
	return spot.NewSpotWebSocketClient(options...)
}

func NewFuturesClient(options ...client.Options) *futures.FuturesClient {
	return futures.NewFuturesClient(options...)
}

func NewFuturesWebSocketClient(options ...client.Options) *futures.FuturesWebSocketClient {
	return futures.NewFuturesWebSocketClient(options...)
}
