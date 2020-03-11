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
	g.printSetBase(len(values), typeName, cfg.setDelimiter, cfg.strictSet)
	g.Printf("\n")
	g.buildSetNoOpOrderChangeDetect(flags, typeName)
	g.Printf("\n")
	g.printSetFlags(flags, values, typeName)
	g.Printf("\n")
}

// Arguments to format are:
//  [1]: type name
//  [2]: base type name
//  [3]: base type bit count
//  [4]: delimiter
//  [5]: parse error handling
//  [6]: parse error description
const setType = `
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

// IsA%[1]sSet returns "true" if the value is a set of values listed in the %[1]s enum definition. "false" otherwise
func (i %[1]sSet) IsA%[1]sSet() bool {
	return i&_%[1]sSetBitmask == i
}
`

func (g *Generator) printSetBase(valueCount int, typeName string, delimiter string, strict bool) {
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
		// TODO: handle valueCount > 64, consider using math/big
		panic("set-of-enums for more than 64 enum value defined is not supported")
	}
	errorHandling := "continue"
	errorHandlingDescription := "Ignores"
	if strict {
		errorHandling = fmt.Sprintf("return 0, fmt.Errorf(\"%%s does not belong to %s values\", s)", typeName)
		errorHandlingDescription = "Returns error if there are"
	}

	g.Printf(setType, typeName, intType, intType[4:], delimiter, errorHandling, errorHandlingDescription)
}

// Arguments to format are:
//  [1]: type name
const setFlagHeader = `
// %[1]sSet flags
const (
`

func (g *Generator) printSetFlags(flags []string, values []Value, typeName string) {
	g.Printf(setFlagHeader, typeName)
	g.Printf("\t%[1]s %[2]sSet = 1 << iota\n\t%[3]s\n)\n\n", flags[0], typeName, strings.Join(flags[1:], "\n\t"))

	g.Printf("var _%[1]sToFlagMap = map[%[1]s]%[1]sSet{\n", typeName)
	for i, flagName := range flags {
		g.Printf("\t%[1]s: %[2]s,\n", values[i].originalName, flagName)
	}
	g.Printf("}\n\n")

	g.Printf("const _%[1]sSetBitmask = %[2]s\n", typeName, strings.Join(flags, " |\n\t"))
}

// buildSetNoOpOrderChangeDetect try to let the compiler and the user know if the order of the ENUMS have changed.
func (g *Generator) buildSetNoOpOrderChangeDetect(flags []string, typeName string) {
	g.Printf("\n")

	g.Printf(`
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	`)
	g.Printf("func _%[1]sSetNoOp (){ ", typeName)
	g.Printf("\n var x [1]struct{}\n")
	for i, flagName := range flags {
		g.Printf("\t_ = x[%s-(%d)]\n", flagName, 1<<i)
	}
	g.Printf("}\n\n")
}
