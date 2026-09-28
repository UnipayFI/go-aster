package common

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

func TestJSONSupport(t *testing.T) {
	if errJSONSupport != nil {
		t.Fatal(errJSONSupport)
	}
}

type unixTimes struct {
	S  time.Time `json:"s,format:unix"`
	Ms time.Time `json:"ms,format:unixmilli"`
	Us time.Time `json:"us,format:unixmicro"`
	Ns time.Time `json:"ns,format:unixnano"`
}

// TestDecodeUnixMatchesStandard checks that every unix unit decodes exactly as
// the standard library's `format` option (and the time package) define it:
// the same instant, in UTC.
func TestDecodeUnixMatchesStandard(t *testing.T) {
	ref := time.Unix(1790597484, 538123456)
	tests := []struct {
		in   string
		want unixTimes
	}{
		{
			in: `{"s":1790597484,"ms":1790597484538,"us":1790597484538123,"ns":1790597484538123456}`,
			want: unixTimes{
				S:  time.Unix(ref.Unix(), 0),
				Ms: time.UnixMilli(ref.UnixMilli()),
				Us: time.UnixMicro(ref.UnixMicro()),
				Ns: ref,
			},
		},
		{
			in: `{"s":1790597484.538123456,"ms":1790597484538.123456,"us":1790597484538123.456}`,
			want: unixTimes{
				S:  ref,
				Ms: ref,
				Us: ref,
			},
		},
		{
			in: `{"s":0,"ms":0,"us":0,"ns":0}`,
			want: unixTimes{
				S:  time.Unix(0, 0),
				Ms: time.Unix(0, 0),
				Us: time.Unix(0, 0),
				Ns: time.Unix(0, 0),
			},
		},
		{
			in: `{"s":-1,"ms":-1,"us":-1,"ns":-1}`,
			want: unixTimes{
				S:  time.Unix(-1, 0),
				Ms: time.UnixMilli(-1),
				Us: time.UnixMicro(-1),
				Ns: time.Unix(0, -1),
			},
		},
		{
			in: `{"ms":253402300799999}`,
			want: unixTimes{
				Ms: time.Date(9999, 12, 31, 23, 59, 59, 999e6, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		var got, std unixTimes
		if err := JSONUnmarshal([]byte(tt.in), &got); err != nil {
			t.Fatalf("JSONUnmarshal(%s): %v", tt.in, err)
		}
		if err := json.Unmarshal([]byte(tt.in), &std, jsonexp.ExperimentalSupportFormatTag(true)); err != nil {
			t.Fatalf("json.Unmarshal(%s): %v", tt.in, err)
		}
		for i, pair := range [][3]time.Time{
			{got.S, std.S, tt.want.S},
			{got.Ms, std.Ms, tt.want.Ms},
			{got.Us, std.Us, tt.want.Us},
			{got.Ns, std.Ns, tt.want.Ns},
		} {
			g, s, w := pair[0], pair[1], pair[2]
			if w.IsZero() {
				if !g.IsZero() {
					t.Errorf("%s: field %d = %v, want zero time", tt.in, i, g)
				}
				continue
			}
			if !g.Equal(w) || !g.Equal(s) || g.Location() != time.UTC {
				t.Errorf("%s: field %d = %v (%v), standard %v, want %v in UTC", tt.in, i, g, g.Location(), s, w)
			}
		}
	}
}

type timeFields struct {
	Ms time.Time  `json:"ms,format:unixmilli"`
	P  *time.Time `json:"p,format:unixmilli"`
}

// TestDecodeUnsetValues pins down how "not set" reaches the SDK: an explicit 0
// is the Unix epoch, while a missing key or null leaves a value field at the
// zero time and a pointer field nil.
func TestDecodeUnsetValues(t *testing.T) {
	epoch := time.Unix(0, 0)
	tests := []struct {
		in    string
		ms    time.Time
		p     *time.Time
		isNil bool
	}{
		{in: `{}`, isNil: true},
		{in: `{"ms":null,"p":null}`, isNil: true},
		{in: `{"ms":0,"p":0}`, ms: epoch, p: &epoch},
		{in: `{"ms":1790597484538,"p":1790597484538}`, ms: time.UnixMilli(1790597484538), p: ptr(time.UnixMilli(1790597484538))},
	}
	for _, tt := range tests {
		var got timeFields
		if err := JSONUnmarshal([]byte(tt.in), &got); err != nil {
			t.Fatalf("JSONUnmarshal(%s): %v", tt.in, err)
		}
		if !got.Ms.Equal(tt.ms) || got.Ms.IsZero() != tt.ms.IsZero() {
			t.Errorf("%s: ms = %v, want %v", tt.in, got.Ms, tt.ms)
		}
		switch {
		case tt.isNil && got.P != nil:
			t.Errorf("%s: p = %v, want nil", tt.in, *got.P)
		case !tt.isNil && (got.P == nil || !got.P.Equal(*tt.p) || got.P.Location() != time.UTC):
			t.Errorf("%s: p = %v, want %v in UTC", tt.in, got.P, *tt.p)
		}
	}
}

func ptr[T any](v T) *T { return &v }

// TestDecodeRejectsMalformed checks that values Aster never sends (quoted
// numbers, empty strings, exponents, JSON of the wrong kind) are errors rather
// than silently misdated.
func TestDecodeRejectsMalformed(t *testing.T) {
	for _, in := range []string{
		`{"ms":"1790597484538"}`,
		`{"ms":""}`,
		`{"ms":"0"}`,
		`{"ms":1.790597484538e12}`,
		`{"ms":true}`,
		`{"ms":{}}`,
		`{"p":"1790597484538"}`,
	} {
		var got timeFields
		if err := JSONUnmarshal([]byte(in), &got); err == nil {
			t.Errorf("JSONUnmarshal(%s) = %+v, want error", in, got)
		}
	}
}

// TestStringOption checks the standard `,string` + format combination: the
// value must be quoted, and a bare number is rejected.
func TestStringOption(t *testing.T) {
	var v struct {
		Ms time.Time `json:"ms,string,format:unixmilli"`
	}
	if err := JSONUnmarshal([]byte(`{"ms":"1790597484538"}`), &v); err != nil || !v.Ms.Equal(time.UnixMilli(1790597484538)) {
		t.Errorf("quoted: got %v, %v", v.Ms, err)
	}
	if err := JSONUnmarshal([]byte(`{"ms":1790597484538}`), &v); err == nil {
		t.Error("bare number with ,string: want error")
	}
	b, err := JSONMarshal(v)
	if err != nil || string(b) != `{"ms":"1790597484538"}` {
		t.Errorf("JSONMarshal = %s, %v", b, err)
	}
}

// TestEncodeUnix checks encoding in each unit: sub-unit precision is kept as
// a fraction (as the standard library does), not truncated.
func TestEncodeUnix(t *testing.T) {
	v := unixTimes{
		S:  time.Unix(1790597484, 0),
		Ms: time.Unix(1790597484, 538000000),
		Us: time.Unix(1790597484, 538123456),
		Ns: time.Unix(1790597484, 538123456).In(time.FixedZone("UTC+8", 8*3600)),
	}
	b, err := JSONMarshal(v)
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"s":1790597484,"ms":1790597484538,"us":1790597484538123.456,"ns":1790597484538123456}`
	if string(b) != want {
		t.Errorf("JSONMarshal = %s, want %s", b, want)
	}
	var back unixTimes
	if err := JSONUnmarshal(b, &back); err != nil || !back.Us.Equal(v.Us) || !back.Ns.Equal(v.Ns) {
		t.Errorf("round trip = %+v, %v", back, err)
	}

	b, err = JSONMarshal(timeFields{})
	if err != nil || string(b) != `{"ms":-62135596800000,"p":null}` {
		t.Errorf("JSONMarshal(zero) = %s, %v", b, err)
	}
	var zero timeFields
	if err := JSONUnmarshal(b, &zero); err != nil || !zero.Ms.IsZero() || zero.P != nil {
		t.Errorf("zero round trip = %+v, %v", zero, err)
	}
}

// TestLayoutFormats checks layout-based formats and the RFC 3339 default for
// a time.Time without a format option.
func TestLayoutFormats(t *testing.T) {
	type layouts struct {
		D time.Time  `json:"d,format:DateOnly"`
		C time.Time  `json:"c,format:'2006-01-02T15:04'"`
		N *time.Time `json:"n,format:RFC3339Nano"`
		R time.Time  `json:"r"`
	}
	const in = `{"d":"2026-09-28","c":"2026-09-28T12:34","n":"2026-09-28T12:34:56.123456789Z","r":"2026-09-28T12:34:56.789+08:00"}`
	var v layouts
	if err := JSONUnmarshal([]byte(in), &v); err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC); !v.D.Equal(want) {
		t.Errorf("d = %v, want %v", v.D, want)
	}
	if want := time.Date(2026, 9, 28, 12, 34, 0, 0, time.UTC); !v.C.Equal(want) {
		t.Errorf("c = %v, want %v", v.C, want)
	}
	if want := time.Date(2026, 9, 28, 12, 34, 56, 123456789, time.UTC); v.N == nil || !v.N.Equal(want) {
		t.Errorf("n = %v, want %v", v.N, want)
	}
	if want := time.Date(2026, 9, 28, 4, 34, 56, 789e6, time.UTC); !v.R.Equal(want) {
		t.Errorf("r = %v, want %v", v.R, want)
	}
	b, err := JSONMarshal(v)
	if err != nil || string(b) != in {
		t.Errorf("JSONMarshal = %s, %v, want %s", b, err, in)
	}
	if err := JSONUnmarshal([]byte(`{"d":1790597484538}`), &v); err == nil {
		t.Error("number into DateOnly: want error")
	}
	if err := JSONUnmarshal([]byte(`{"r":1790597484538}`), &v); err == nil {
		t.Error("number into untagged time.Time: want error")
	}
}

// TestOtherStandardFormats checks that the `format` options of other types are
// honoured too, since the codec enables format tags as a whole.
func TestOtherStandardFormats(t *testing.T) {
	type other struct {
		D  time.Duration  `json:"d,format:units"`
		Ds time.Duration  `json:"ds,format:sec"`
		Di time.Duration  `json:"di,format:iso8601"`
		B  []byte         `json:"b,format:hex"`
		B6 []byte         `json:"b6,format:base64url"`
		F  float64        `json:"f,format:nonfinite"`
		M  map[string]int `json:"m,format:emitnull"`
		S  []int          `json:"s,format:emitnull"`
		Me map[string]int `json:"me,format:emitempty"`
		Se []int          `json:"se,format:emitempty"`
	}
	v := other{D: 90 * time.Second, Ds: 1500 * time.Millisecond, Di: time.Hour, B: []byte{1, 0xab}, B6: []byte{0xfb, 0xff}, F: math.Inf(-1)}
	const want = `{"d":"1m30s","ds":1.5,"di":"PT1H","b":"01ab","b6":"-_8=","f":"-Infinity","m":null,"s":null,"me":{},"se":[]}`
	b, err := JSONMarshal(v)
	if err != nil || string(b) != want {
		t.Fatalf("JSONMarshal = %s, %v, want %s", b, err, want)
	}
	var back other
	if err := JSONUnmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	if back.D != v.D || back.Ds != v.Ds || back.Di != v.Di || string(back.B) != string(v.B) || string(back.B6) != string(v.B6) || !math.IsInf(back.F, -1) {
		t.Errorf("round trip = %+v, want %+v", back, v)
	}
	if err := JSONUnmarshal([]byte(`{"f":"NaN"}`), &back); err != nil || !math.IsNaN(back.F) {
		t.Errorf("NaN: %v, %v", back.F, err)
	}
}

// TestDecimal checks decimal.Decimal, which keeps its own JSON methods: quoted
// or bare on input, quoted on output.
func TestDecimal(t *testing.T) {
	var v struct {
		D decimal.Decimal `json:"d"`
	}
	for in, want := range map[string]string{`{"d":"0.00012345"}`: "0.00012345", `{"d":83417.24338768}`: "83417.24338768", `{"d":"-1e-8"}`: "-0.00000001"} {
		if err := JSONUnmarshal([]byte(in), &v); err != nil || v.D.String() != want {
			t.Errorf("JSONUnmarshal(%s) = %v, %v, want %s", in, v.D, err, want)
		}
	}
	v.D = decimal.RequireFromString("1.50")
	if b, err := JSONMarshal(v); err != nil || string(b) != `{"d":"1.5"}` {
		t.Errorf("JSONMarshal = %s, %v", b, err)
	}
}

// TestFormatTagNeedsCodec documents why SDK types must go through
// JSONMarshal/JSONUnmarshal: plain encoding/json/v2 and encoding/json reject
// the `format` tag option.
func TestFormatTagNeedsCodec(t *testing.T) {
	var v timeFields
	errs := []error{json.Unmarshal([]byte(`{"ms":1}`), &v), jsonv1.Unmarshal([]byte(`{"ms":1}`), &v)}
	_, err := json.Marshal(v)
	errs = append(errs, err)
	_, err = jsonv1.Marshal(v)
	errs = append(errs, err)
	for i, err := range errs {
		if err == nil || !strings.Contains(err.Error(), "format") {
			t.Errorf("case %d: err = %v, want unsupported format tag error", i, err)
		}
	}
}

func TestConcurrentUse(t *testing.T) {
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 200 {
				ms := int64(1790597484538 + i*1000 + j)
				in := fmt.Sprintf(`{"ms":%d,"p":%d}`, ms, ms+1)
				var v timeFields
				if err := JSONUnmarshal([]byte(in), &v); err != nil || v.Ms.UnixMilli() != ms || v.P.UnixMilli() != ms+1 {
					t.Errorf("JSONUnmarshal(%s) = %+v, %v", in, v, err)
					return
				}
				if b, err := JSONMarshal(v); err != nil || string(b) != in {
					t.Errorf("JSONMarshal = %s, %v, want %s", b, err, in)
					return
				}
			}
		}()
	}
	wg.Wait()
}
