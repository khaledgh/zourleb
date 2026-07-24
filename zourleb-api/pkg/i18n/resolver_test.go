package i18n

import "testing"

type tr struct {
	locale string
	val    string
}

func (t tr) GetLocale() string { return t.locale }

func TestResolveFallbackChain(t *testing.T) {
	items := []tr{{"en", "Hello"}, {"ar", "مرحبا"}, {"fr", "Bonjour"}}

	// requested locale present
	if got, ok := Resolve(items, "fr", "en"); !ok || got.val != "Bonjour" {
		t.Errorf("requested fr = %q ok=%v", got.val, ok)
	}
	// requested missing → default
	if got, ok := Resolve(items, "de", "en"); !ok || got.val != "Hello" {
		t.Errorf("fallback to default = %q ok=%v", got.val, ok)
	}
	// requested & default missing → first available
	if got, ok := Resolve(items, "de", "es"); !ok || got.val != "Hello" {
		t.Errorf("fallback to first = %q ok=%v", got.val, ok)
	}
	// empty slice
	if _, ok := Resolve([]tr{}, "en", "en"); ok {
		t.Errorf("empty slice should return ok=false")
	}
}
