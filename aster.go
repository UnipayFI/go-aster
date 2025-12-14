package aster

import (
	"github.com/UnipayFI/go-aster/client"
	"github.com/UnipayFI/go-aster/futures"
	"github.com/UnipayFI/go-aster/spot"
)

func NewSpotClient(options ...client.Options) *spot.SpotClient {
	return spot.NewSpotClient(options...)
}

func NewFuturesClient(options ...client.Options) *futures.FuturesClient {
	return futures.NewFuturesClient(options...)
}
