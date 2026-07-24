// Package i18n resolves translated content fields with a fallback chain.
package i18n

// Translatable is any sibling translation row exposing its locale.
type Translatable interface {
	GetLocale() string
}

// Resolve picks the best translation from a slice using the chain:
// requested locale → default locale → first available. Returns the zero value
// index (-1) when the slice is empty.
func Resolve[T Translatable](items []T, requested, def string) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	var byDefault, first T
	haveDefault, haveFirst := false, false
	for _, it := range items {
		if it.GetLocale() == requested {
			return it, true
		}
		if it.GetLocale() == def && !haveDefault {
			byDefault, haveDefault = it, true
		}
		if !haveFirst {
			first, haveFirst = it, true
		}
	}
	if haveDefault {
		return byDefault, true
	}
	return first, haveFirst
}
