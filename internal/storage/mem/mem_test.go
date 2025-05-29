package mem_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/inmore/gopaste/internal/model"
	"github.com/inmore/gopaste/internal/storage/mem"
)

func TestStore_SaveLoad(t *testing.T) {
	t.Parallel()

	st := mem.New()

	p := &model.Paste{
		ID:        "id1",
		Content:   "check value in mem-store",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := st.Save(p); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := st.Load("id1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got.Content != p.Content {
		t.Errorf("content mismatch: want %q, got %q", p.Content, got.Content)
	}
}

func TestStore_DeleteExpired(t *testing.T) {
	t.Parallel()

	st := mem.New()

	expired := &model.Paste{
		ID:        "old",
		Content:   "value1",
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	live := &model.Paste{
		ID:        "live",
		Content:   "value2",
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}

	_ = st.Save(expired)
	_ = st.Save(live)

	deleted, err := st.DeleteExpired()
	if err != nil {
		t.Fatalf("DeleteExpired error: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 record deleted, got %d", deleted)
	}

	if _, err := st.Load("old"); err == nil {
		t.Fatal("expired paste should have been removed")
	}
	if _, err := st.Load("live"); err != nil {
		t.Fatalf("live paste disappeared: %v", err)
	}
}

func BenchmarkStore_SaveLoad(b *testing.B) {
	st := mem.New()

	base := &model.Paste{
		Content:   "bench-data",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := "id" + strconv.Itoa(i)
		p := *base
		p.ID = id

		_ = st.Save(&p)
		_, _ = st.Load(id)
	}
}
