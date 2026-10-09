package tools

import "testing"

func TestCalc(t *testing.T) {
	ok := map[string]string{
		"2*(3+4)^2": "98", "-2^2": "-4", "2^3^2": "512", "10/4": "2.5",
		"7%3": "1", " 1 + 2 * 3 ": "7", "-(1+2)": "-3", "0.1+0.2": "0.3",
	}
	for in, want := range ok {
		if got, err := Calc(in); err != nil || got != want {
			t.Errorf("Calc(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "1/0", "(1+2", "1+", "abc", "1 2", "2^", "5%0"} {
		if got, err := Calc(in); err == nil {
			t.Errorf("Calc(%q) = %q; want error", in, got)
		}
	}
}
