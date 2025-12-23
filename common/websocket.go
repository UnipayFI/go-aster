package common

import "strings"

// WsSymbolStreamSegment returns a combined stream segment like `btcusdt@aggTrade`.
// The symbol is lowercased, and non-empty extraParts are appended as `@<part>`.
func WsSymbolStreamSegment(symbol, stream string, extraParts ...string) string {
	segment := strings.ToLower(symbol) + "@" + stream
	for _, part := range extraParts {
		if part == "" {
			continue
		}
		segment += "@" + part
	}
	return segment
}

// WsAllStreamSegment returns a combined stream segment like `!ticker@arr`.
// Non-empty extraParts are appended as `@<part>`.
func WsAllStreamSegment(stream string, extraParts ...string) string {
	segment := "!" + stream
	for _, part := range extraParts {
		if part == "" {
			continue
		}
		segment += "@" + part
	}
	return segment
}

// WsSymbolStream returns a WebSocket endpoint like `/ws/btcusdt@aggTrade`.
func WsSymbolStream(symbol, stream string, extraParts ...string) string {
	return WEBSOCKET_STREAM_SEPARATOR + WsSymbolStreamSegment(symbol, stream, extraParts...)
}

// WsAllStream returns a WebSocket endpoint like `/ws/!ticker@arr`.
func WsAllStream(stream string, extraParts ...string) string {
	return WEBSOCKET_STREAM_SEPARATOR + WsAllStreamSegment(stream, extraParts...)
}

// WsCombinedStreamsEndpoint builds a combined streams endpoint like `/stream?streams=s1/s2/...`.
// Each segment must not include the `/ws/` prefix.
func WsCombinedStreamsEndpoint(segments ...string) string {
	return WEBSOCKET_COMBINED_SEPARATOR + strings.Join(segments, "/")
}
