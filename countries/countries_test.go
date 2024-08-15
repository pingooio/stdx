package countries_test

import (
	"testing"

	"github.com/pingooio/stdx/countries"
)

func TestGetMap(t *testing.T) {
	expectedNumberOfCountries := 249

	err := countries.Init()
	if err != nil {
		t.Error(err)
		return
	}
	countries := countries.GetMap()

	if len(countries) != expectedNumberOfCountries {
		t.Errorf("Invalid number of countries. Got %d, expected: %d", len(countries), expectedNumberOfCountries)
	}
}

func TestGetList(t *testing.T) {
	expectedNumberOfCountries := 249

	err := countries.Init()
	if err != nil {
		t.Error(err)
		return
	}
	countries := countries.GetList()

	if len(countries) != expectedNumberOfCountries {
		t.Errorf("Invalid number of countries. Got %d, expected: %d", len(countries), expectedNumberOfCountries)
	}
}
