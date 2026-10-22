package main

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

const testModule = "github.com/vixartech/example"

func TestLocate(t *testing.T) {
	t.Parallel()

	tests := append(validLocateCases(), invalidLocateCases()...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, actualApp, err := locate(testModule, test.importPath)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("locate(%q) error = %v, expected %v", test.importPath, err, test.expectedErr)
			}

			if actualApp != test.expectedApp {
				t.Fatalf("locate(%q) application = %t, expected %t", test.importPath, actualApp, test.expectedApp)
			}

			if actual != test.expected {
				t.Errorf("locate(%q) = %#v, expected %#v", test.importPath, actual, test.expected)
			}
		})
	}
}

func TestPermits(t *testing.T) {
	t.Parallel()

	tests := append(corePermissionCases(), boundaryPermissionCases()...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := permits(test.source, test.target)
			if actual != test.expected {
				t.Errorf("permits(%#v, %#v) = %t, expected %t", test.source, test.target, actual, test.expected)
			}
		})
	}
}

func TestValidatePackages(t *testing.T) {
	t.Parallel()

	tests := validationCases()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := validatePackages(testModule, test.packages)
			if !slices.Equal(actual, test.expected) {
				t.Errorf("validatePackages() = %q, expected %q", actual, test.expected)
			}
		})
	}
}

func TestDecodePackages(t *testing.T) {
	t.Parallel()

	tests := decodePackageCases()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := decodePackages(strings.NewReader(test.input))
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("decodePackages() error = %v, expected %v", err, test.expectedErr)
			}

			if !slices.EqualFunc(actual, test.expected, equalPackageInfo) {
				t.Errorf("decodePackages() = %#v, expected %#v", actual, test.expected)
			}
		})
	}
}

type decodePackageCase struct {
	name        string
	input       string
	expected    []packageInfo
	expectedErr error
}

func decodePackageCases() []decodePackageCase {
	return []decodePackageCase{
		{
			name:  "go list package",
			input: `{"ImportPath":"example/domain","Imports":["errors"]}`,
			expected: []packageInfo{
				{ImportPath: "example/domain", Imports: []string{"errors"}},
			},
		},
		{
			name:  "package without imports",
			input: `{"ImportPath":"example/domain"}`,
			expected: []packageInfo{
				{ImportPath: "example/domain"},
			},
		},
		{name: "missing import path", input: `{}`, expectedErr: errMissingPath},
	}
}

func equalPackageInfo(left packageInfo, right packageInfo) bool {
	return left.ImportPath == right.ImportPath && slices.Equal(left.Imports, right.Imports)
}

type locateCase struct {
	name        string
	importPath  string
	expected    packageLocation
	expectedApp bool
	expectedErr error
}

func validLocateCases() []locateCase {
	return []locateCase{
		{
			name:        "existing service query package",
			importPath:  appPath("service/query/problems"),
			expected:    location(layerService, ""),
			expectedApp: true,
		},
		{
			name:        "query named transport helper",
			importPath:  appPath("ports/http/query"),
			expected:    location(layerPorts, "http"),
			expectedApp: true,
		},
		{
			name:        "nested query package",
			importPath:  appPath("query/problems/window"),
			expected:    location("query", ""),
			expectedApp: true,
		},
		{
			name:        "service package",
			importPath:  appPath("service/command"),
			expected:    location(layerService, ""),
			expectedApp: true,
		},
		{
			name:        "adapter owner",
			importPath:  appPath("adapters/postgres/gen"),
			expected:    location(layerAdapters, "postgres"),
			expectedApp: true,
		},
		{name: "external package", importPath: "net/http"},
	}
}

func invalidLocateCases() []locateCase {
	return []locateCase{
		invalidLocateCase("application root", testModule+"/internal/app/example"),
		invalidLocateCase("unknown layer", appPath("runtime")),
		invalidLocateCase("ownerless adapter", appPath("adapters")),
	}
}

func TestQueryImports(t *testing.T) {
	t.Parallel()

	for _, test := range queryImportCases() {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			current := packageInfo{ImportPath: appPath(test.source), Imports: []string{test.target}}

			actual := validatePackages(testModule, []packageInfo{current})
			if (len(actual) == 0) != test.allowed {
				t.Errorf("%s imports %s: violations = %q, expected allowed = %t",
					current.ImportPath, test.target, actual, test.allowed)
			}
		})
	}
}

type queryImportCase struct {
	name, source, target string
	allowed              bool
}

func queryImportCases() []queryImportCase {
	return []queryImportCase{
		{"query standard library", "query", "context", true},
		{"query domain", "query", appPath("domain"), true},
		{"nested query domain", "query/problems/window", appPath("domain/value"), true},
		{"query sibling", "query/problems", appPath("query/source"), true},
		{"port query", "ports/http", appPath("query"), true},
		{"adapter query", "adapters/postgres", appPath("query"), true},
		{"query service", "query", appPath("service"), false},
		{"nested query command", "query/problems/window", appPath("service/command"), false},
		{"query port", "query", appPath("ports/http"), false},
		{"query adapter", "query/problems", appPath("adapters/postgres"), false},
		{"domain query", "domain", appPath("query"), false},
		{"service query", "service", appPath("query"), false},
		{"query lookalike", "query", appPath("querying"), false},
		{"service adapter", "service", appPath("adapters/postgres"), false},
		{"query other app", "query", testModule + "/internal/app/other/domain", false},
		{"query other query", "query/problems", testModule + "/internal/app/other/query", false},
		{"port other query", "ports/http", testModule + "/internal/app/other/query", false},
	}
}

func invalidLocateCase(name string, importPath string) locateCase {
	return locateCase{
		name:        name,
		importPath:  importPath,
		expectedApp: true,
		expectedErr: errInvalidLocation,
	}
}

type permissionCase struct {
	name     string
	source   packageLocation
	target   packageLocation
	expected bool
}

func corePermissionCases() []permissionCase {
	return []permissionCase{
		permission("domain to domain", layerDomain, "", layerDomain, "", true),
		permission("domain to service", layerDomain, "", layerService, "", false),
		permission("service to domain", layerService, "", layerDomain, "", true),
		permission("service to adapter", layerService, "", layerAdapters, "postgres", false),
	}
}

func boundaryPermissionCases() []permissionCase {
	return []permissionCase{
		permission("port to service", layerPorts, "http", layerService, "", true),
		permission("same port", layerPorts, "http", layerPorts, "http", true),
		permission("sibling port", layerPorts, "http", layerPorts, "grpc", false),
		permission("adapter to service", layerAdapters, "postgres", layerService, "", true),
		permission("same adapter", layerAdapters, "postgres", layerAdapters, "postgres", true),
		permission("sibling adapter", layerAdapters, "postgres", layerAdapters, "redis", false),
		{
			name:   "other service",
			source: location(layerService, ""),
			target: packageLocation{service: "other", layer: layerDomain},
		},
	}
}

func permission(
	name string,
	sourceLayer string,
	sourceOwner string,
	targetLayer string,
	targetOwner string,
	expected bool,
) permissionCase {
	return permissionCase{
		name:     name,
		source:   location(sourceLayer, sourceOwner),
		target:   location(targetLayer, targetOwner),
		expected: expected,
	}
}

type validationCase struct {
	name     string
	packages []packageInfo
	expected []string
}

func validationCases() []validationCase {
	return []validationCase{
		{
			name: "valid service dependencies",
			packages: []packageInfo{
				packageWithImport("service", "domain"),
			},
		},
		{
			name: "sibling adapter dependency",
			packages: []packageInfo{
				packageWithImport("adapters/postgres", "adapters/redis"),
			},
			expected: []string{siblingAdapterViolation()},
		},
	}
}

func packageWithImport(source string, target string) packageInfo {
	return packageInfo{
		ImportPath: appPath(source),
		Imports:    []string{appPath(target)},
	}
}

func siblingAdapterViolation() string {
	return "- " + appPath("adapters/postgres") + " (adapters) must not import " +
		appPath("adapters/redis") + " (adapters)"
}

func location(layer string, owner string) packageLocation {
	return packageLocation{service: "example", layer: layer, owner: owner}
}

func appPath(suffix string) string {
	return testModule + "/internal/app/example/" + suffix
}
