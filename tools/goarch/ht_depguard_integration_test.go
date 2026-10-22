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
