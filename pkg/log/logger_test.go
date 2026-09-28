package log

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

// TestDefaultLoggerFormats checks that the printf-style methods format their
// arguments into the message instead of handing them to slog as attributes,
// and that disabled levels log nothing.
func TestDefaultLoggerFormats(t *testing.T) {
	var buf bytes.Buffer
	l := &defaultLogger{Logger: slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))}

	l.Errorf("request %s failed: %s", "https://sapi.asterdex.com/api/v3/depth", errors.New("boom"))
	l.Warnf("retry %d of %d", 1, 3)
	l.Infof("offset=%dms", 42)
	l.Debugf("response: %s", "hidden")

	out := buf.String()
	for _, want := range []string{
		`level=ERROR msg="request https://sapi.asterdex.com/api/v3/depth failed: boom"`,
		`level=WARN msg="retry 1 of 3"`,
		`level=INFO msg="offset=42ms"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %s:\n%s", want, out)
		}
	}
	if strings.Contains(out, "%") || strings.Contains(out, "BADKEY") || strings.Contains(out, "hidden") {
		t.Errorf("unformatted, attribute-style or disabled output:\n%s", out)
	}
}
