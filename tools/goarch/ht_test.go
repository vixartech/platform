package main

import (
	"errors"
	"testing"
)

func htApps() []string {
	return []string{"spidertower", "spiderssr", "spiderspa", "htworker", "htindexer", "htapi", "html"}
}

func TestHTLocations(t *testing.T) {
	t.Parallel()

	for _, app := range htApps() {
		t.Run(app, func(t *testing.T) {
			t.Parallel()

			testHTLocations(t, app)
		})
	}
}

func testHTLocations(t *testing.T, app string) {
	t.Helper()

	for _, test := range htLocationCases() {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			path := testModule + "/internal/app/" + app + "/" + test.importPath

			actual, local, err := locate(testModule, path)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("locate(%q) error = %v, expected %v", path, err, test.expectedErr)
			}

			if !local || actual.layer != test.expected.layer {
				t.Errorf("locate(%q) = %#v, local = %t, expected layer %q", path, actual, local, test.expected.layer)
			}
		})
	}
}

func htLocationCases() []locateCase {
	return []locateCase{
		{name: "command", importPath: "service/command", expected: location(layerCommand, "")},
		{name: "nested command", importPath: "service/command/work", expected: location(layerCommand, "")},
		{name: "query", importPath: "service/query", expected: location(layerQuery, "")},
		{name: "nested query", importPath: "service/query/graph", expected: location(layerQuery, "")},
		{name: "missing operation", importPath: "service", expectedErr: errInvalidLocation},
		{name: "unknown operation", importPath: "service/work", expectedErr: errInvalidLocation},
		{name: "top level query", importPath: "query/graph", expectedErr: errInvalidLocation},
	}
}

func TestHTImports(t *testing.T) {
	t.Parallel()

	for _, app := range htApps() {
		t.Run(app, func(t *testing.T) {
			t.Parallel()

			for _, test := range htImportCases() {
				t.Run(test.name, func(t *testing.T) {
					t.Parallel()

					base := testModule + "/internal/app/" + app + "/"
					current := packageInfo{ImportPath: base + test.source, Imports: []string{base + test.target}}

					violations := validatePackages(testModule, []packageInfo{current})
					if (len(violations) == 0) != test.allowed {
						t.Errorf("%s imports %s: violations = %q, expected allowed = %t",
							test.source, test.target, violations, test.allowed)
					}
				})
			}
		})
	}
}

func htImportCases() []queryImportCase {
	return []queryImportCase{
		{"command domain", "service/command", "domain", true},
		{"query domain", "service/query", "domain", true},
		{"command sibling", "service/command/work", "service/command/policy", true},
		{"query sibling", "service/query/graph", "service/query/model", true},
		{"command query", "service/command", "service/query", false},
		{"query command", "service/query", "service/command", false},
		{"nested command query", "service/command/work", "service/query/graph", false},
		{"nested query command", "service/query/graph", "service/command/work", false},
		{"domain command", "domain", "service/command", false},
		{"domain query", "domain", "service/query", false},
		{"port command", "ports/http", "service/command", true},
		{"port query", "ports/http", "service/query", true},
		{"adapter command", "adapters/postgres", "service/command", true},
		{"adapter query", "adapters/postgres", "service/query", true},
		{"command adapter", "service/command", "adapters/postgres", false},
		{"query port", "service/query", "ports/http", false},
		{"command lookalike", "service/command", "service/commanding", false},
		{"query lookalike", "service/query", "service/querying", false},
	}
}
