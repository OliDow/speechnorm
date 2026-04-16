// Package speechnorm converts digits in text into locale-appropriate words
// suitable for TTS input.
//
// The package is self-contained and intentionally has no dependencies
// outside the Go standard library. It may be extracted to its own repo
// in future.
package speechnorm

import "sync"

// Converter converts integer values to word representations in a given
// locale.
type Converter interface {
	// ToWords returns the cardinal word form of n.
	ToWords(n int64) string
	// ToOrdinalWords returns the ordinal word form of n. Locales without
	// an ordinal regex trigger may return the cardinal form.
	ToOrdinalWords(n int64) string
}

var (
	registryMu sync.RWMutex
	registry   = map[string]Converter{}
)

// Register binds a locale code (e.g. "en", "fr") to a Converter. It is
// intended to be called from each locale file's init(). Panics if the
// locale has already been registered — duplicate registration is always
// a programming error.
func Register(locale string, c Converter) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[locale]; exists {
		panic("speechnorm: locale already registered: " + locale)
	}
	registry[locale] = c
}

// lookup returns the Converter for locale, or (nil, false) if none is
// registered.
func lookup(locale string) (Converter, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	c, ok := registry[locale]
	return c, ok
}

// snapshotRegistry and restoreRegistry are test helpers that let tests
// mutate the registry without affecting other tests.
func snapshotRegistry() map[string]Converter {
	registryMu.RLock()
	defer registryMu.RUnlock()
	snap := make(map[string]Converter, len(registry))
	for k, v := range registry {
		snap[k] = v
	}
	return snap
}

func restoreRegistry(snap map[string]Converter) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = snap
}
