package speechnorm

import (
	"testing"
)

type fakeConv struct{ name string }

func (f fakeConv) ToWords(_ int64) string        { return f.name }
func (f fakeConv) ToOrdinalWords(_ int64) string { return f.name + "-ord" }

func TestRegisterAndLookup(t *testing.T) {
	isolatedRegistry(t)

	Register("xx", fakeConv{name: "xx"})
	got, ok := lookup("xx")
	if !ok {
		t.Fatalf("lookup(xx) returned ok=false, want true")
	}
	if got.ToWords(0) != "xx" {
		t.Errorf("got %q, want %q", got.ToWords(0), "xx")
	}
}

func TestLookupUnknownReturnsFalse(t *testing.T) {
	isolatedRegistry(t)
	if _, ok := lookup("zz"); ok {
		t.Errorf("lookup(zz) returned ok=true, want false")
	}
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	isolatedRegistry(t)
	Register("yy", fakeConv{name: "yy"})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on duplicate register")
		}
	}()
	Register("yy", fakeConv{name: "yy2"})
}

func isolatedRegistry(t *testing.T) {
	t.Helper()
	snap := snapshotRegistry()
	registryMu.Lock()
	registry = map[string]Converter{}
	registryMu.Unlock()
	t.Cleanup(func() { restoreRegistry(snap) })
}
