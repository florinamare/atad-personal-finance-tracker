package money

import (
	"encoding/json"
	"testing"
)

func TestParseValid(t *testing.T) {
	cases := map[string]Amount{
		"-12.50": -1250,
		"0.10":   10,
		"100.00": 10000,
		"-0.00":  0,
		"0.01":   1,
	}
	for in, want := range cases {
		got, err := Parse(in)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("Parse(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	bad := []string{"", "12", "12.5", "12.500", "1,50", "abc", "-", ".50",
		"+1.00", "1.-5", "--1.00", " 1.00", "1.00 ", "12345678901234.00"}
	for _, in := range bad {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) should fail", in)
		}
	}
}

func TestString(t *testing.T) {
	cases := map[Amount]string{-1250: "-12.50", 10: "0.10", 0: "0.00", -5: "-0.05", 100000: "1000.00"}
	for in, want := range cases {
		if got := in.String(); got != want {
			t.Errorf("Amount(%d).String() = %q, want %q", in, got, want)
		}
	}
}

// The test from the project description: -0.10 + -0.20 must be exactly "-0.30".
func TestSumIsExact(t *testing.T) {
	a, _ := Parse("-0.10")
	b, _ := Parse("-0.20")
	if got := (a + b).String(); got != "-0.30" {
		t.Errorf("sum = %q, want \"-0.30\"", got)
	}
}

func TestJSON(t *testing.T) {
	var a Amount
	if err := json.Unmarshal([]byte(`"-0.10"`), &a); err != nil || a != -10 {
		t.Errorf("unmarshal string: got %d, err %v", a, err)
	}
	if err := json.Unmarshal([]byte(`-0.10`), &a); err == nil {
		t.Error("unmarshal of a JSON number should fail")
	}
	out, _ := json.Marshal(Amount(-30))
	if string(out) != `"-0.30"` {
		t.Errorf("marshal = %s, want \"-0.30\"", out)
	}
}