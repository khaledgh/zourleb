// Package slug builds URL-safe slugs with a short uniqueness suffix.
package slug

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Make normalizes s to a slug. Non-ASCII (e.g. Arabic) input that reduces to
// empty falls back to a random token so the slug is always usable.
func Make(s string) string {
	out := nonAlnum.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "-")
	out = strings.Trim(out, "-")
	if out == "" {
		out = "item"
	}
	return out
}

// MakeUnique appends a short random suffix to guarantee uniqueness.
func MakeUnique(s string) string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return Make(s) + "-" + hex.EncodeToString(b)
}
