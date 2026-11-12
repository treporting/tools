package create

import (
	"strings"
)

// marks holds, per SPDX identifier, the sentences that its licence text holds and no other of the
// list does. They are compared with the text, its spacing ignored.
var marks = []struct {
	id   string
	all  []string
	none []string
}{
	{id: "Apache-2.0", all: []string{"Apache License", "Version 2.0"}},
	{id: "MIT", all: []string{"Permission is hereby granted, free of charge"}},
	{id: "BSD-3-Clause", all: []string{"Redistribution and use in source and binary forms", "Neither the name"}},
	{id: "BSD-2-Clause", all: []string{"Redistribution and use in source and binary forms"}, none: []string{"Neither the name"}},
	{id: "ISC", all: []string{"Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee"}},
	{id: "MPL-2.0", all: []string{"Mozilla Public License", "2.0"}},
}

// Licence answers the licence of a licence text by its SPDX identifier, several joined by AND, and
// unknown where it recognizes none.
//
// @implements SR_43@1.1.1
func Licence(text string) string {
	flat := strings.Join(strings.Fields(text), " ")
	var ids []string
	for _, m := range marks {
		if holdsAll(flat, m.all) && !holdsAny(flat, m.none) {
			ids = append(ids, m.id)
		}
	}
	if len(ids) == 0 {
		return Unknown
	}
	return strings.Join(ids, " AND ")
}

// Given answers a licence a package names itself, such as the field license of package-lock.json,
// where it has the form of an SPDX expression - identifiers joined by AND, OR or WITH - and unknown
// where it names none or names it in words, such as "2-clause BSD-like license".
//
// @implements SR_43@1.1.1
func Given(licence string) string {
	licence = strings.TrimSpace(licence)
	if !spdxForm(licence) {
		return Unknown
	}
	return licence
}

// spdxForm tells whether an expression is identifiers joined by AND, OR or WITH, in parentheses or
// not. It checks the form only, not whether SPDX lists the identifier.
func spdxForm(expression string) bool {
	tokens := strings.Fields(strings.NewReplacer("(", " ", ")", " ").Replace(expression))
	if len(tokens)%2 == 0 {
		return false
	}
	for i, token := range tokens {
		operator := token == "AND" || token == "OR" || token == "WITH"
		if (i%2 == 1) != operator || (!operator && !identifier(token)) {
			return false
		}
	}
	return true
}

func identifier(token string) bool {
	for _, r := range token {
		allowed := r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || strings.ContainsRune(".-+:", r)
		if !allowed {
			return false
		}
	}
	return token != ""
}

func holdsAll(text string, sentences []string) bool {
	for _, s := range sentences {
		if !strings.Contains(text, s) {
			return false
		}
	}
	return true
}

func holdsAny(text string, sentences []string) bool {
	for _, s := range sentences {
		if strings.Contains(text, s) {
			return true
		}
	}
	return false
}
