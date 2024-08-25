package countries_test

import (
	"testing"

	"github.com/pingooio/stdx/countries"
)

func TestGetMap(t *testing.T) {
	expectedNumberOfCountries := 249

	countries := countries.All()

	if len(countries) != expectedNumberOfCountries {
		t.Errorf("Invalid number of countries. Got %d, expected: %d", len(countries), expectedNumberOfCountries)
	}
}

func TestGetCountry(t *testing.T) {
	tests := []struct {
		code string
		name string
	}{
		{"FR", "France"},
		{"XX", "Unknown"},
	}

	for _, test := range tests {
		countryName, _ := countries.Name(test.code)
		if countryName != test.name {
			t.Errorf("Code: %s -> Expected: %s | Got: %s", test.code, test.name, countryName)
		}
	}
}
