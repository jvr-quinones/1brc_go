package station_test

import (
	. "1brc_go/station"
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

func TestParseLineFloat(t *testing.T) {
	for _, test := range readlineTestArray {
		_, actualFloat, actualErr := ParseLineFloat(test.text)
		if actualErr != test.expectErr || test.expectFloat != actualFloat {
			t.Errorf("Error parsing %q: Got %f %v instead of %f %v",
				test.text, actualFloat, actualErr, test.expectFloat, test.expectErr)
		}
	}
}

func TestParseLineInt(t *testing.T) {
	for _, test := range readlineTestArray {
		_, actualInt, actualErr := ParseLineInt(test.text)
		if actualErr != test.expectErr || test.expectInt != actualInt {
			t.Errorf("Error parsing %q: Got %d %v instead of %d %v",
				test.text, actualInt, actualErr, test.expectInt, test.expectErr)
		}
	}
}

func BenchmarkParseLineFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, _, err := ParseLineFloat("Las Palmas de Gran Canaria;31.4"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLineInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, _, err := ParseLineInt("Yakutsk;-347"); err != nil {
			b.Fatal(err)
		}
	}
}
