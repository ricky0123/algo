package main

import (
	"fmt"
	"strings"
)

type Alias struct {
	Identifier string
	Command    string
	SubAliases []*Alias
}

func PrettyPrint(aliases []*Alias) string {
	var b strings.Builder

	shebang := fmt.Sprintf("#!/usr/bin/env bash\n\n")
	b.WriteString(shebang)

	for _, alias := range aliases {
		b.WriteString(alias.PrettyPrint())
	}

	return b.String()
}

func (a Alias) PrettyPrint() string {
	var b strings.Builder

	preamble := fmt.Sprintf("\n%v() {\n", a.Identifier)
	b.WriteString(preamble)

	body := a.prettyPrintSubAlias(1, "")
	b.WriteString(body)

	b.WriteString("}\n")
	return b.String()
}

func (a Alias) prettyPrintSubAlias(level int, commandPrefix string) string {
	var b strings.Builder
	indent := strings.Repeat("\t", level)

	switch {
	case commandPrefix == "" && a.Command != "":
		commandPrefix = a.Command
	case commandPrefix != "" && a.Command != "":
		commandPrefix += " " + a.Command
	}

	b.WriteString(
		fmt.Sprintf("%vcase \"$1\" in\n", indent),
	)
	for _, subAlias := range a.SubAliases {
		b.WriteString(
			fmt.Sprintf("%v\t%v)\n", indent, subAlias.Identifier),
		)
		b.WriteString(
			fmt.Sprintf("%v\t\tshift\n", indent),
		)
		subAliasSource := subAlias.prettyPrintSubAlias(level+2, commandPrefix)
		b.WriteString(subAliasSource)
		b.WriteString(
			fmt.Sprintf("%v\t\t;;\n", indent),
		)
	}
	b.WriteString(
		fmt.Sprintf("%v\t*)\n", indent),
	)
	b.WriteString(
		fmt.Sprintf("%v\t\t%v \"$@\"\n", indent, commandPrefix),
	)
	b.WriteString(
		fmt.Sprintf("%v\t\treturn #?\n", indent),
	)
	b.WriteString(
		fmt.Sprintf("%v\t\t;;\n", indent),
	)
	b.WriteString(
		fmt.Sprintf("%vesac\n", indent),
	)

	return b.String()
}
