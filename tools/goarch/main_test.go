package main

import (
	"errors"
	"slices"
	"testing"
)

const testModule = "github.com/vixartech/example"

func TestLocate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		importPath  string
		expected    packageLocation
		expectedApp bool
		expectedErr error
	}{
		{
			name:        "service package",
			importPath:  testModule + "/internal/app/example/service/command",
			expected:    packageLocation{service: "example", layer: "service", owner: "command"},
			expectedApp: true,
		},
		{
			name:        "adapter owner",
			importPath:  testModule + "/internal/app/example/adapters/postgres/gen",
			expected:    packageLocation{service: "example", layer: "adapters", owner: "postgres"},
			expectedApp: true,
		},
		{name: "external package", importPath: "net/http"},
		{
			name:        "application root",
			importPath:  testModule + "/internal/app/example",
			expectedApp: true,
			expectedErr: errInvalidLocation,
		},
		{
			name:        "unknown layer",
			importPath:  testModule + "/internal/app/example/runtime",
			expectedApp: true,
			expectedErr: errInvalidLocation,
		},
		{
			name:        "ownerless adapter",
			importPath:  testModule + "/internal/app/example/adapters",
			expectedApp: true,
			expectedErr: errInvalidLocation,
		},
	}

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

	tests := []struct {
		name     string
		source   packageLocation
		target   packageLocation
		expected bool
	}{
		{name: "domain to domain", source: location("domain", ""), target: location("domain", ""), expected: true},
		{name: "domain to service", source: location("domain", ""), target: location("service", "")},
		{name: "service to domain", source: location("service", ""), target: location("domain", ""), expected: true},
		{name: "service to adapter", source: location("service", ""), target: location("adapters", "postgres")},
		{name: "port to service", source: location("ports", "http"), target: location("service", ""), expected: true},
		{name: "same port", source: location("ports", "http"), target: location("ports", "http"), expected: true},
		{name: "sibling port", source: location("ports", "http"), target: location("ports", "grpc")},
		{name: "adapter to service", source: location("adapters", "postgres"), target: location("service", ""), expected: true},
		{
			name:     "same adapter",
			source:   location("adapters", "postgres"),
			target:   location("adapters", "postgres"),
			expected: true,
		},
		{
			name:   "sibling adapter",
			source: location("adapters", "postgres"),
			target: location("adapters", "redis"),
		},
		{
			name:   "other service",
			source: location("service", ""),
			target: packageLocation{service: "other", layer: "domain"},
		},
	}

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

	tests := []struct {
		name     string
		packages []packageInfo
		expected []string
	}{
		{
			name: "valid service dependencies",
			packages: []packageInfo{
				{
					ImportPath: testModule + "/internal/app/example/service",
					Imports:    []string{testModule + "/internal/app/example/domain"},
				},
			},
		},
		{
			name: "sibling adapter dependency",
			packages: []packageInfo{
				{
					ImportPath: testModule + "/internal/app/example/adapters/postgres",
					Imports:    []string{testModule + "/internal/app/example/adapters/redis"},
				},
			},
			expected: []string{
				"- github.com/vixartech/example/internal/app/example/adapters/postgres (adapters) must not import " +
					"github.com/vixartech/example/internal/app/example/adapters/redis (adapters)",
			},
		},
	}

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

func location(layer string, owner string) packageLocation {
	return packageLocation{service: "example", layer: layer, owner: owner}
}
