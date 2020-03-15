// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Set test: enumeration of more than 64 values.

package main

import (
	"fmt"
	"math/big"
)

type Country int

const (
	Afghanistan Country = iota
	Albania
	Algeria
	Andorra
	Angola
	AntiguaAndBarbuda
	Argentina
	Armenia
	Australia
	Austria
	Azerbaijan
	TheBahamas
	Bahrain
	Bangladesh
	Barbados
	Belarus
	Belgium
	Belize
	Benin
	Bhutan
	Bolivia
	BosniaAndHerzegovina
	Botswana
	Brazil
	Brunei
	Bulgaria
	BurkinaFaso
	Burundi
	CaboVerde
	Cambodia
	Cameroon
	Canada
	CentralAfricanRepublic
	Chad
	Chile
	China
	Colombia
	Comoros
	DemocraticRepublicOfTheCongo
	RepublicOfTheCongo
	CostaRica
	CôteDIvoire
	Croatia
	Cuba
	Cyprus
	CzechRepublic
	Denmark
	Djibouti
	Dominica
	DominicanRepublic
	EastTimor
	Ecuador
	Egypt
	ElSalvador
	EquatorialGuinea
	Eritrea
	Estonia
	Eswatini
	Ethiopia
	Fiji
	Finland
	France
	Gabon
	TheGambia
	Georgia
	Germany
	Ghana
	Greece
	Grenada
	Guatemala
	Guinea
	GuineaBissau
	Guyana
	Haiti
	Honduras
	Hungary
	Iceland
	India
	Indonesia
	Iran
	Iraq
	Ireland
	Israel
	Italy
	Jamaica
	Japan
	Jordan
	Kazakhstan
	Kenya
	Kiribati
	NorthKorea
	SouthKorea
	Kosovo
	Kuwait
	Kyrgyzstan
	Laos
	Latvia
	Lebanon
	Lesotho
	Liberia
	Libya
	Liechtenstein
	Lithuania
	Luxembourg
	Madagascar
	Malawi
	Malaysia
	Maldives
	Mali
	Malta
	MarshallIslands
	Mauritania
	Mauritius
	Mexico
	FederatedStatesOfMicronesia
	Moldova
	Monaco
	Mongolia
	Montenegro
	Morocco
	Mozambique
	Myanmar
	Namibia
	Nauru
	Nepal
	Netherlands
	NewZealand
	Nicaragua
	Niger
	Nigeria
	NorthMacedonia
	Norway
	Oman
	Pakistan
	Palau
	Panama
	PapuaNewGuinea
	Paraguay
	Peru
	Philippines
	Poland
	Portugal
	Qatar
	Romania
	Russia
	Rwanda
	SaintKittsAndNevis
	SaintLucia
	SaintVincentAndTheGrenadines
	Samoa
	SanMarino
	SaoTomeAndPrincipe
	SaudiArabia
	Senegal
	Serbia
	Seychelles
	SierraLeone
	Singapore
	Slovakia
	Slovenia
	SolomonIslands
	Somalia
	SouthAfrica
	Spain
	SriLanka
	Sudan
	SouthSudan
	Suriname
	Sweden
	Switzerland
	Syria
	Taiwan
	Tajikistan
	Tanzania
	Thailand
	Togo
	Tonga
	TrinidadAndTobago
	Tunisia
	Turkey
	Turkmenistan
	Tuvalu
	Uganda
	Ukraine
	UnitedArabEmirates
	UnitedKingdom
	UnitedStates
	Uruguay
	Uzbekistan
	Vanuatu
	VaticanCity
	Venezuela
	Vietnam
	Yemen
	Zambia
	Zimbabwe
)

func main() {
	ck(CountrySetHungary, "Hungary")
	ck(CountrySetVaticanCity, "VaticanCity")
	ck(CountrySet{big.NewInt(3)}, "Afghanistan Albania")
	ck(CountrySet{big.NewInt(0).Or(CountrySetHungary.Int, CountrySetAustria.Int)}, "Austria Hungary")
	ck(CountrySet{big.NewInt(0).Or(CountrySetUnitedKingdom.Int, CountrySetUnitedStates.Int)}, "UnitedKingdom UnitedStates")
	ckCountrySetElements(CountrySetHungary, []Country{Hungary})
	ckCountrySetElements(
		CountrySet{big.NewInt(0).Or(CountrySetUnitedKingdom.Int, CountrySetUnitedStates.Int)},
		[]Country{UnitedKingdom, UnitedStates},
	)
	ckCountrySetTest("Austria Hungary", "VaticanCity", false)
	ckCountrySetTest("Austria Hungary", "UnitedKingdom UnitedStates", false)
	ckCountrySetTest("Austria Hungary", "UnitedKingdom Austria", false)
	ckCountrySetTest("Austria Hungary", "Austria", true)
	ckCountrySetTest("Austria Hungary", "Hungary", true)
	ckCountrySetTest("Austria Hungary", "Hungary Austria", true)
	ckCountrySetString(CountrySetHungary, "Hungary")
	ckCountrySetString(CountrySetHungary, "hungary")
	ckCountrySetString(CountrySet{big.NewInt(0).Or(CountrySetHungary.Int, CountrySetAustria.Int)}, "Austria Hungary")
	ckCountrySetString(CountrySetHungary, "Hungary Antarctica")
}

var panicPrefix = "set_big.go: "

func ck(country CountrySet, str string) {
	if country.String() != str {
		panic(panicPrefix + str)
	}
}

func ckCountrySetElements(countrySet CountrySet, expectedElems []Country) {
	elems := countrySet.Elements()
	if len(elems) != len(expectedElems) {
		panic(fmt.Sprintf("%sexpected %d elements, got %d", panicPrefix, len(expectedElems), len(elems)))
	}
OUTER:
	for _, expected := range expectedElems {
		for _, got := range elems {
			if got == expected {
				continue OUTER
			}
		}
		panic(fmt.Sprintf("%smissing element %s", panicPrefix, expected))
	}
}

func ckCountrySetTest(cs string, m string, expected bool) {
	countrySet, err := CountrySetString(cs)
	ckNoError(err)
	mask, err := CountrySetString(m)
	ckNoError(err)
	if countrySet.Test(mask) != expected {
		panic(fmt.Sprintf("%stesting '%s', on '%s' is not %v", panicPrefix, mask.String(), countrySet.String(), expected))
	}
}

func ckCountrySetString(country CountrySet, str string) {
	c, err := CountrySetString(str)
	if err != nil {
		panic(panicPrefix + err.Error())
	}
	if c.Cmp(country.Int) != 0 {
		panic(panicPrefix + str)
	}
}

func ckNoError(err error) {
	if err != nil {
		panic(panicPrefix + err.Error())
	}
}
