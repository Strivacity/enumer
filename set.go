package main

import (
	"fmt"
	"strings"
)

func (g *Generator) addSetType(values []Value, typeName string, cfg config) {
	trimPrefixes := strings.Split(cfg.trimPrefix, ",")
	flags := make([]string, len(values))
	for i, v := range values {
		flagName := v.originalName
		for _, prefix := range trimPrefixes {
			flagName = strings.TrimPrefix(flagName, prefix)
		}
		flags[i] = fmt.Sprintf("%[1]sSet%[2]s", typeName, flagName)
	}

	g.Printf("\n")
	g.printSetBase(flags, typeName, cfg.setDelimiter, cfg.strictSet)
	g.Printf("\n")
	g.printSetFlagMap(flags, values, typeName)
	g.Printf("\n")
}

// Arguments to format are:
//  [1]: type name
//  [2]: base type name
//  [3]: base type bit count
//  [4]: delimiter
//  [5]: parse error handling
//  [6]: parse error description
//  [7]: enum value count
const setType = `
// this is necessary because of import "math/big"
var _ = big.NewInt(0)

// %[1]sSet is a combination of %[1]s values
type %[1]sSet %[2]s

func (i %[1]sSet) String() string {
	elements := i.Elements()
	elementStrings := make([]string, len(elements))
	for idx, e := range elements {
		elementStrings[idx] = e.String()
	}
	return strings.Join(elementStrings, "%[4]s")
}

func (i %[1]sSet) Elements() []%[1]s {
	elements := make([]%[1]s, 0, bits.OnesCount%[3]s(%[2]s(i)))
	for _, value := range %[1]sValues() {
		if i.Test(_%[1]sToFlagMap[value]) {
			elements = append(elements, value)
		}
	}
	return elements
}

// Test returns true if all %[1]s values in mask are set, false otherwise.
func (i %[1]sSet) Test(mask %[1]sSet) bool {
	return mask.IsA%[1]sSet() && i&mask == mask
}

// %[1]sSetString retrieves an enum set value from the "%[4]s" delimited list of %[1]s constant string names.
// %[6]s elements in the list which are not part of the %[1]s enum.
func %[1]sSetString(s string) (%[1]sSet, error) {
	i := %[1]sSet(0)
	for _, s := range strings.Split(s, "%[4]s") {
		v, err := %[1]sString(s)
		if err != nil {
			%[5]s
		}
		flag, ok := _%[1]sToFlagMap[v]
		if !ok { // this should not happen if the generator is correct
			%[5]s
		}
		i |= flag
	}
	return i, nil
}

const _%[1]sSetBitmask  = 1<<(%[7]d+1) - 1

// IsA%[1]sSet returns "true" if the value is a set of values listed in the %[1]s enum definition. "false" otherwise
func (i %[1]sSet) IsA%[1]sSet() bool {
	return i&_%[1]sSetBitmask == i
}
`

// Arguments to format are:
//  [1]: type name
//  [2]: base type package
//  [3]: base type name
//  [4]: delimiter
//  [5]: parse error handling
//  [6]: parse error description
//  [7]: enum value count
const setTypeBigInt = `
// this is necessary because of import "math/bits"
var _ = bits.TrailingZeros(0)
// %[1]sSet is a combination of %[1]s values
type %[1]sSet struct{*%[2]s.%[3]s}

func (i %[1]sSet) String() string {
	elements := i.Elements()
	elementStrings := make([]string, len(elements))
	for idx, e := range elements {
		elementStrings[idx] = e.String()
	}
	return strings.Join(elementStrings, "%[4]s")
}

func (i %[1]sSet) Elements() []%[1]s {
	values := %[1]sValues()
	elements := make([]%[1]s, 0, i.BitLen()-int(i.TrailingZeroBits()))
	for idx := 0; idx < i.BitLen(); idx++ {
		if i.Bit(idx) == 1 {
			elements = append(elements, values[idx])
		}
	}
	return elements
}

// Test returns true if all %[1]s values in mask are set, false otherwise.
func (i %[1]sSet) Test(mask %[1]sSet) bool {
	return mask.IsA%[1]sSet() && %[2]s.New%[3]s(0).And(i.%[3]s, mask.%[3]s).Cmp(mask.%[3]s) == 0
}

// %[1]sSetString retrieves an enum set value from the "%[4]s" delimited list of %[1]s constant string names.
// %[6]s elements in the list which are not part of the %[1]s enum.
func %[1]sSetString(s string) (%[1]sSet, error) {
	i := %[1]sSet{%[2]s.New%[3]s(0)}
	for _, s := range strings.Split(s, "%[4]s") {
		v, err := %[1]sString(s)
		if err != nil {
			%[5]s
		}
		flag, ok := _%[1]sToFlagMap[v]
		if !ok { // this should not happen if the generator is correct
			%[5]s
		}
		i.Or(i.%[3]s, flag.%[3]s)
	}
	return i, nil
}

var _%[1]sSetBitmask  = %[2]s.New%[3]s(0).Sub(%[2]s.New%[3]s(0).SetBit(%[2]s.New%[3]s(0), %[7]d+1, 1), %[2]s.New%[3]s(1))

// IsA%[1]sSet returns "true" if the value is a set of values listed in the %[1]s enum definition. "false" otherwise
func (i %[1]sSet) IsA%[1]sSet() bool {
	return i.CmpAbs(_%[1]sSetBitmask) <= 0
}
`

func (g *Generator) printSetBase(flags []string, typeName string, delimiter string, strict bool) {
	const bigInt = "big.Int"
	valueCount := len(flags)
	intType := ""
	switch {
	case valueCount <= 8:
		intType = "uint8"
	case valueCount <= 16:
		intType = "uint16"
	case valueCount <= 32:
		intType = "uint32"
	case valueCount <= 64:
		intType = "uint64"
	default:
		intType = bigInt
	}
	errorHandling := "continue"
	errorHandlingDescription := "Ignores"
	if strict {
		errorHandling = fmt.Sprintf("return 0, fmt.Errorf(\"%%s does not belong to %s values\", s)", typeName)
		errorHandlingDescription = "Returns error if there are"
	}

	if intType == bigInt {
		externalType := strings.Split(intType, ".")
		g.Printf(setTypeBigInt, typeName, externalType[0], externalType[1], delimiter, errorHandling, errorHandlingDescription, valueCount)
		g.Printf(`
		// %[1]sSet flags
		var (
		`, typeName)
		for idx := range flags {
			g.Printf("\t%[1]s = %[2]sSet{big.NewInt(0).SetBit(big.NewInt(0), %[3]d, 1)}\n", flags[idx], typeName, idx)
		}
		g.Printf(")\n\n")
	} else {
		g.Printf(setType, typeName, intType, intType[4:], delimiter, errorHandling, errorHandlingDescription, valueCount)
		g.Printf(`
		// %[1]sSet flags
		const (
		`, typeName)
		g.Printf("\t%[1]s %[2]sSet = 1 << iota\n\t%[3]s\n)\n\n", flags[0], typeName, strings.Join(flags[1:], "\n\t"))
		g.Printf("\n")
		g.buildSetNoOpOrderChangeDetect(flags, typeName)
	}
}

func (g *Generator) printSetFlagMap(flags []string, values []Value, typeName string) {
	g.Printf("var _%[1]sToFlagMap = map[%[1]s]%[1]sSet{\n", typeName)
	for i, flagName := range flags {
		g.Printf("\t%[1]s: %[2]s,\n", values[i].originalName, flagName)
	}
	g.Printf("}\n\n")
}

// buildSetNoOpOrderChangeDetect try to let the compiler and the user know if the order of the ENUMS have changed.
func (g *Generator) buildSetNoOpOrderChangeDetect(flags []string, typeName string) {
	g.Printf(`
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	func _%[1]sSetNoOp (){
		var x [1]struct{}
	`, typeName)
	for i, flagName := range flags {
		g.Printf("\t_ = x[%s-(%d)]\n", flagName, 1<<i)
	}
	g.Printf("}\n\n")
}

// buildSetNoOpOrderChangeDetectBigInt try to let the compiler and the user know if the order of the ENUMS have changed.
func (g *Generator) buildSetNoOpOrderChangeDetectBigInt(flags []string, typeName string) {
	g.Printf(`
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	func _%[1]sSetNoOp (){
		var x [1]struct{}
	`, typeName)
	for i, flagName := range flags {
		g.Printf("\t_ = x[big.NewInt(0).Sub(%s, big.NewInt(0).SetBit(big.NewInt(0), %d, 1)).Int64()]\n", flagName, i)
	}
	g.Printf("}\n\n")
}
