package common

import (
	"encoding/json/v2"
	"fmt"
	"time"

	jsonexp "github.com/go-json-experiment/json"
)

// Aster sends timestamps as bare JSON numbers. Every time.Time field declares
// the unit it arrives in with the standard `format` tag option (e.g.
// `json:"updateTime,format:unixmilli"`), which Go 1.27's encoding/json/v2 only
// honours when ExperimentalSupportFormatTag is set, so SDK types are encoded
// and decoded through JSONMarshal/JSONUnmarshal.
var (
	jsonOptions json.Options

	// errJSONSupport is non-nil when the running Go release no longer
	// provides the `format` tag support this codec relies on.
	errJSONSupport = initJSON()
)

// JSONMarshal marshals v, honouring `format` tag options.
func JSONMarshal(v any) ([]byte, error) {
	if errJSONSupport != nil {
		return nil, errJSONSupport
	}
	return json.Marshal(v, jsonOptions)
}

// JSONUnmarshal unmarshals data into v, honouring `format` tag options.
func JSONUnmarshal(data []byte, v any) error {
	if errJSONSupport != nil {
		return errJSONSupport
	}
	return json.Unmarshal(data, v, jsonOptions)
}

// initJSON builds the codec options and round-trips a probe through them, so
// a Go release that drops the experimental `format` tag support (by panicking
// on the unknown option or by ignoring it) fails loudly instead of silently
// misdating fields.
func initJSON() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("aster: encoding/json/v2 format tag support unavailable: %v", r)
		}
	}()
	jsonOptions = jsonexp.ExperimentalSupportFormatTag(true)

	const payload = `{"t":1750034396998123}`
	var probe struct {
		T time.Time `json:"t,format:unixmicro"`
	}
	if err := json.Unmarshal([]byte(payload), &probe, jsonOptions); err != nil {
		return fmt.Errorf("aster: encoding/json/v2 format tag support unavailable: %w", err)
	}
	if got := probe.T.UnixMicro(); got != 1750034396998123 {
		return fmt.Errorf("aster: encoding/json/v2 format tag probe decoded %d, want 1750034396998123", got)
	}
	if out, err := json.Marshal(probe, jsonOptions); err != nil || string(out) != payload {
		return fmt.Errorf("aster: encoding/json/v2 format tag probe encoded %s (%v), want %s", out, err, payload)
	}
	return nil
}
