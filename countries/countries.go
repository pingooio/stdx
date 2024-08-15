package countries

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
)

// sources
// https://gist.github.com/keeguon/2310008
// see also ttps://restcountries.eu/rest/v2/all

//go:embed countries.json
var Bytes []byte

const (
	Unknown     = "Unknown"
	CodeUnknown = "XX"
)

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var countriesMap map[string]Country
var countriesList []Country

var (
	ErrCountryNotFound = errors.New("Country not found")
)

func Init() (err error) {
	countriesList = []Country{}

	err = json.Unmarshal(Bytes, &countriesList)
	if err != nil {
		return fmt.Errorf("countries: error decoding JSON: %w", err)
	}

	countriesMap = map[string]Country{}
	for _, country := range countriesList {
		if _, exists := countriesMap[country.Code]; exists {
			err = fmt.Errorf("countries: ducplicate country for country code: %s", country.Code)
			return
		}

		countriesMap[country.Code] = country
	}

	return
}

func GetMap() map[string]Country {
	return countriesMap
}

func GetList() []Country {
	return countriesList
}

func Name(countryCode string) (countryName string, err error) {
	countryName = Unknown

	country, exists := countriesMap[countryCode]
	if !exists {
		err = ErrCountryNotFound
		return
	}

	countryName = country.Name
	return
}
