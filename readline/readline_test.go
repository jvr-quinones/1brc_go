package readline_test

import (
	. "1brc_go/readline"
	"testing"
)

type readlineTest struct {
	text        string
	expectFloat float64
	expectInt   int64
	expectErr   error
}

var readlineTestArray = []readlineTest{
	{"Test01:1.0", 0.0, 0, ErrNoSep},
	{"Test02:-2.0", 0.0, 0, ErrNoSep},
	{"Test03;3.0", 3.0, 30, nil},
	{"Test04;-4.0", -4.0, -40, nil},
}

func TestReadAsFloat(t *testing.T) {
	for _, test := range readlineTestArray {
		_, actualFloat, actualErr := ReadAsFloat(test.text)
		if actualErr != test.expectErr || test.expectFloat != actualFloat {
			t.Errorf("Error parsing %q: Got %f %v instead of %f %v",
				test.text, actualFloat, actualErr, test.expectFloat, test.expectErr)
		}
	}
}

func TestReadAsInt(t *testing.T) {
	for _, test := range readlineTestArray {
		_, actualInt, actualErr := ReadAsInt(test.text)
		if actualErr != test.expectErr || test.expectInt != actualInt {
			t.Errorf("Error parsing %q: Got %d %v instead of %d %v",
				test.text, actualInt, actualErr, test.expectInt, test.expectErr)
		}
	}
}

func BenchmarkReadFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, _, err := ReadAsFloat("Las Palmas de Gran Canaria;31.4"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, _, err := ReadAsInt("Yakutsk;-347"); err != nil {
			b.Fatal(err)
		}
	}
}
