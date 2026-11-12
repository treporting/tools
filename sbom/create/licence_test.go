package create

import "testing"

const (
	mitText    = "Permission is hereby granted, free of charge, to any person obtaining a copy"
	apacheText = "Apache License\n                           Version 2.0, January 2004"
	bsd3Text   = "Redistribution and use in source and binary forms, with or without\nmodification, are permitted ...\n   * Neither the name of Google Inc. nor the names of its contributors"
	bsd2Text   = "Redistribution and use in source and binary forms, with or without modification, are permitted"
)

// @implements TC_113@1.1.1
func TestLicenceByItsSPDXIdentifier(t *testing.T) {
	for text, want := range map[string]string{
		mitText:    "MIT",
		apacheText: "Apache-2.0",
		bsd3Text:   "BSD-3-Clause",
		bsd2Text:   "BSD-2-Clause",
		"Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted": "ISC",
	} {
		if got := Licence(text); got != want {
			t.Errorf("Licence(%.30q) = %q, want %q", text, got, want)
		}
	}
}

// @implements TC_114@1.1.1
func TestSeveralLicencesJoinedByAnd(t *testing.T) {
	// as the LICENSE of gopkg.in/yaml.v3 has it
	got := Licence("This project is covered by two different licenses: MIT and Apache.\n" + mitText + "\n" + apacheText)
	if got != "Apache-2.0 AND MIT" {
		t.Errorf("got %q", got)
	}
}

// @implements TC_115@1.1.1
func TestLicenceNotRecognizedIsUnknown(t *testing.T) {
	for _, text := range []string{"", "All rights reserved.", "Do what you want."} {
		if got := Licence(text); got != Unknown {
			t.Errorf("Licence(%q) = %q, want unknown", text, got)
		}
	}
}

// @verifies SR_43@1.1.1
func TestGivenLicenceInTheFormOfSPDX(t *testing.T) {
	for given, want := range map[string]string{
		"MIT":                 "MIT",
		" GPL-2.0-only ":      "GPL-2.0-only",
		"(MIT OR Apache-2.0)": "(MIT OR Apache-2.0)",
		"MPL-2.0 AND MIT":     "MPL-2.0 AND MIT",
		"GPL-2.0-or-later WITH Classpath-exception-2.0": "GPL-2.0-or-later WITH Classpath-exception-2.0",
		"":                          Unknown,
		"2-clause BSD-like license": Unknown,
		"MIT AND":                   Unknown,
		"AND MIT":                   Unknown,
		"MIT/X11":                   Unknown,
	} {
		if got := Given(given); got != want {
			t.Errorf("Given(%q) = %q, want %q", given, got, want)
		}
	}
}
