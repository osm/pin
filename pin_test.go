package pin

import (
	"testing"
)

// testGender contains gender based test cases.
var testGender = []struct {
	pin    string
	isMale bool
	err    string
}{
	{"19901121-8774", true, ""},
	{"19131221-7324", false, ""},
}

// TestGender makes sure that the gender function works.
func TestGender(t *testing.T) {
	for _, tc := range testGender {
		m, _ := IsMale(tc.pin)
		if tc.isMale == true && m != true {
			t.Errorf("%s should be interpreted as a male, but it wasn't", tc.pin)
		}

		f, _ := IsFemale(tc.pin)
		if tc.isMale == false && f != true {
			t.Errorf("%s should be interpreted as a female, but it wasn't", tc.pin)
		}
	}
}

// testValidity contains test cases for the IsValid function.
var testValidity = []struct {
	pin string
	err bool
}{
	{"9011218774", true},
	{"901121-8774", true},
	{"901121+8774", true},
	{"19901121-8774", false},
	{"19901121+8774", false},
}

// TestValidty tests that IsValid works.
func TestValidty(t *testing.T) {
	for _, tc := range testValidity {
		v, err := IsValid(tc.pin)
		if tc.err && err == nil {
			t.Errorf("%s is not a valid personal identity number, this should have failed", tc.pin)
		}
		if tc.err == false && v != true {
			t.Errorf("%s is a valid personal identity number, this should not have failed", tc.pin)
		}
	}
}

// testGenerateFromDate contains test cases for the GenerateFromDate function.
var testGenerateFromDate = []struct {
	date string
}{
	{"19901121"},
	{"19840707"},
}

// TestGenerateFromDate tests that GenerateFromDate function works.
func TestGenerateFromDate(t *testing.T) {
	for _, tc := range testGenerateFromDate {
		pin, err := GenerateFromDate(tc.date)
		if err != nil {
			t.Errorf("%s should return an error, the error was %v", tc.date, err)
		}
		_, err = IsValid(pin)
		if err != nil {
			t.Errorf("%s is not a valid personal identity number, the error was: %v", tc.date, err)
		}
	}
}

// TestGenerate tests that the Generate function works.
func TestGenerate(t *testing.T) {
	pin, err := Generate()
	if err != nil {
		t.Errorf("%v", err)
	}
	_, err = IsValid(pin)
	if err != nil {
		t.Errorf("%v", err)
	}
}
