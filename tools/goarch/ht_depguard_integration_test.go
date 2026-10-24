//go:build integration

package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

//nolint:paralleltest // Run real linter processes sequentially to bound CPU and memory use.
func TestHTDepguard(t *testing.T) {
	linter := os.Getenv("GOARCH_TEST_LINTER")
	if linter == "" {
		t.Skip("set GOARCH_TEST_LINTER to the pinned golangci-lint binary")
	}

	config, err := filepath.Abs("../../configs/go/golangci.yml")
	if err != nil {
		t.Fatal(err)
	}

	binary := guardGoBinary(t)
	env := guardEnvironment(binary)

	for _, app := range htApps() {
		t.Run(app, func(t *testing.T) {
			root := guardRoot(t, app)
			writeGuardFixture(t, root, "xnet/go.mod", "module golang.org/x/net\n\ngo 1.26.5\n")
			expected := htGuardFixtures(t, root, app)
			vendorGuardFixture(t, binary, root, env)
			output := runGuardLinter(t, linter, config, root, env)
			assertGuardOutput(t, expected, output)
		})
	}
}

func htGuardFixtures(t *testing.T, root, app string) map[string]bool {
	t.Helper()

	expected := make(map[string]bool)

	for index, test := range htImportCases() {
		for _, nested := range []bool{false, true} {
			base := "internal/app/" + app + "/"
			suffix := strconv.Itoa(index)

			source := base + test.source + "/check" + suffix + ".go"
			if nested {
				source = base + test.source + "/case" + suffix + "/check.go"
			}

			target := base + test.target + "/target" + suffix
			imported := "github.com/vixartech/" + app + "/" + target
			writeGuardFixture(t, root, source, "package fixture\nimport _ "+strconv.Quote(imported)+"\n")
			writeGuardFixture(t, root, target+"/empty.go", "package fixture\n")
			expected[filepath.Join(root, filepath.FromSlash(source))] = !test.allowed
		}
	}

	return expected
}

//nolint:paralleltest,maintidx // One real linter run compares allowed image import and forbidden boundary escapes.
func TestSupportImageDepguard(t *testing.T) {
	linter := os.Getenv("GOARCH_TEST_LINTER")
	if linter == "" {
		t.Skip("set GOARCH_TEST_LINTER to the pinned golangci-lint binary")
	}

	config, err := filepath.Abs("../../configs/go/golangci.yml")
	if err != nil {
		t.Fatal(err)
	}

	binary := guardGoBinary(t)
	env := guardEnvironment(binary)
	root := guardRoot(t, "support")
	module := "module github.com/vixartech/support\n\ngo 1.26.5\n\n" +
		"require golang.org/x/image v0.0.0\nreplace golang.org/x/image => ./ximage\n"
	writeGuardFixture(t, root, "go.mod", module)
	writeGuardFixture(t, root, "ximage/go.mod", "module golang.org/x/image\n\ngo 1.26.5\n")
	writeGuardFixture(t, root, "ximage/webp/empty.go", "package webp\n")
	writeGuardFixture(t, root, "ximage/vp8/empty.go", "package vp8\n")

	expected := make(map[string]bool)

	for _, test := range []struct {
		path, imported string
		denied         bool
	}{
		{path: "ports/http", imported: "webp"},
		{path: "domain", imported: "webp", denied: true},
		{path: "ports/other", imported: "vp8", denied: true},
	} {
		source := "internal/app/support/" + test.path + "/image.go"
		body := "package fixture\nimport _ " + strconv.Quote("golang.org/x/image/"+test.imported) + "\n"
		writeGuardFixture(t, root, source, body)
		expected[filepath.Join(root, filepath.FromSlash(source))] = test.denied
	}

	vendorGuardFixture(t, binary, root, env)
	assertGuardOutput(t, expected, runGuardLinter(t, linter, config, root, env))
}
