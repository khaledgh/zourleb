package phone

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"+961 3 123456", "+9613123456", true},
		{"0096170123456", "+96170123456", true},
		{"03 123 456", "+9613123456", true},
		{"70123456", "+96170123456", true},
		{"+96170123456", "+96170123456", true},
		{"96170123456", "+96170123456", true},
		{"", "", false},
		{"123", "", false},
		{"+1 202 555 0100", "", false}, // non-Lebanese
	}
	for _, c := range cases {
		got, err := Normalize(c.in)
		if c.ok && err != nil {
			t.Errorf("Normalize(%q) unexpected error: %v", c.in, err)
			continue
		}
		if !c.ok && err == nil {
			t.Errorf("Normalize(%q) expected error, got %q", c.in, got)
			continue
		}
		if c.ok && got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMask(t *testing.T) {
	if got := Mask("+96170123456"); got != "+96170***456" {
		t.Errorf("Mask = %q", got)
	}
}
