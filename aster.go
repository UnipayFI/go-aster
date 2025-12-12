package aster

import "github.com/UnipayFI/go-aster/spot"

func NewSpotClient(options ...spot.Options) *spot.SpotClient {
	return spot.NewSpotClient(options...)
}
