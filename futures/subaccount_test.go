package futures

import (
	"context"
	"testing"
)

// TestMigrateUserAssetsEmptyResponse checks the empty response the docs
// describe for a source account with nothing to migrate, and a batchId.
func TestMigrateUserAssetsEmptyResponse(t *testing.T) {
	for body, want := range map[string]string{"": "", `{"batchId":"b1"}`: "b1"} {
		got, err := newTestClient(t, body).NewMigrateUserAssetsService(testUserAddress, 1, "0xsig").Do(context.Background())
		if err != nil || got == nil || got.BatchID != want {
			t.Errorf("body %q: got %+v, %v, want batchId %q", body, got, err, want)
		}
	}
	if _, err := newTestClient(t, `{"batchId":1}`).NewMigrateUserAssetsService(testUserAddress, 1, "0xsig").Do(context.Background()); err == nil {
		t.Error("numeric batchId: want error")
	}
}
