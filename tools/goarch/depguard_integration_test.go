//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// The real linter verifies shared YAML matching, including arbitrary checkout names.
func TestQueryDepguard(t *testing.T) {
	linter := os.Getenv("GOARCH_TEST_LINTER")
	if linter == "" {
		t.Skip("set GOARCH_TEST_LINTER to the pinned golangci-lint binary")
	}

	config, err := filepath.Abs("../../configs/go/golangci.yml")
	if err != nil {
		t.Fatal(err)
	}

	for _, registered := range []struct {
		app   string
		ports bool
	}{
		{"gost", true}, {"dozor", true}, {"id", true}, {"mironenko", true},
		{"signal", true}, {"relay", false}, {"srez", true},
	} {
		t.Run(registered.app, func(t *testing.T) {
			app := registered.app
			root := filepath.Join(t.TempDir(), "arbitrary-checkout")
			module := "github.com/vixartech/" + app
			writeGuardFixture(t, root, "go.mod", "module "+module+"\n\ngo 1.26.5\n"+
				"\nrequire golang.org/x/net v0.0.0\nreplace golang.org/x/net => ./xnet\n")
			root, err = filepath.EvalSymlinks(root)
			if err != nil {
				t.Fatal(err)
			}
			expected := make(map[string]bool)
			for index, test := range queryImportCases() {
				if strings.HasPrefix(test.source, "ports/") && !registered.ports {
					continue
				}
				source := "internal/app/" + app + "/" + test.source + "/case" + strconv.Itoa(index) + "/check.go"
				target := strings.ReplaceAll(test.target, testModule, module)
				target = strings.ReplaceAll(target, "/internal/app/example/", "/internal/app/"+app+"/")
				writeGuardFixture(t, root, source, "package fixture\nimport _ "+strconv.Quote(target)+"\n")
				if strings.HasPrefix(target, module+"/") {
					writeGuardFixture(t, root, strings.TrimPrefix(target, module+"/")+"/empty.go", "package fixture\n")
				}
				expected[filepath.Join(root, filepath.FromSlash(source))] = !test.allowed
			}
			// Direct files must match too; distinct targets keep the fixture acyclic.
			for _, direct := range []queryImportCase{
				{source: "query", target: "ports/grpc"},
				{source: "domain", target: "query"},
				{source: "service", target: "query"},
				{source: "ports/http", target: "query", allowed: true},
				{source: "adapters/postgres", target: "query", allowed: true},
			} {
				if strings.HasPrefix(direct.source, "ports/") && !registered.ports {
					continue
				}
				file := "internal/app/" + app + "/" + direct.source + "/direct.go"
				target := module + "/internal/app/" + app + "/" + direct.target
				writeGuardFixture(t, root, file, "package fixture\nimport _ "+strconv.Quote(target)+"\n")
				writeGuardFixture(t, root, strings.TrimPrefix(target, module+"/")+"/empty.go", "package fixture\n")
				expected[filepath.Join(root, filepath.FromSlash(file))] = !direct.allowed
			}
			fixture := module + "/internal/test/fixture"
			writeGuardFixture(t, root, "internal/test/fixture/empty.go", "package fixture\n")
			for _, file := range []string{"fixture.go", "fixture_test.go"} {
				path := "internal/app/" + app + "/service/" + file
				writeGuardFixture(t, root, path, "package fixture\nimport _ "+strconv.Quote(fixture)+"\n")
				expected[filepath.Join(root, filepath.FromSlash(path))] = file == "fixture.go" || app != "srez"
			}
			writeGuardFixture(t, root, "xnet/go.mod", "module golang.org/x/net\n\ngo 1.26.5\n")
			for index, imported := range []string{"golang.org/x/net/idna", "golang.org/x/net/idna/subpkg", "golang.org/x/net/http2"} {
				file := "internal/app/" + app + "/domain/external" + strconv.Itoa(index) + ".go"
				writeGuardFixture(t, root, file, "package fixture\nimport _ "+strconv.Quote(imported)+"\n")
				writeGuardFixture(t, root, "xnet/"+strings.TrimPrefix(imported, "golang.org/x/net/")+"/empty.go", "package fixture\n")
				expected[filepath.Join(root, filepath.FromSlash(file))] = imported != "golang.org/x/net/idna"
			}
			fixtureEnv := append(os.Environ(), "GO111MODULE=on", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off",
				"PATH="+filepath.Join(runtime.GOROOT(), "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
			vendor := exec.CommandContext(t.Context(), filepath.Join(runtime.GOROOT(), "bin", "go"), "mod", "vendor")
			vendor.Dir = root
			vendor.Env = fixtureEnv
			if output, err := vendor.CombinedOutput(); err != nil {
				t.Fatalf("vendor local guard fixture: %v, %s", err, output)
			}

			command := exec.CommandContext(t.Context(), linter, "run", "--config", config,
				"--enable-only", "depguard", "--output.json.path=stdout", "--output.text.path=", "./...")
			command.Dir = root
			command.Env = fixtureEnv
			var stderr bytes.Buffer
			command.Stderr = &stderr
			output, runErr := command.Output()
			var exitErr *exec.ExitError
			if runErr != nil && (!errors.As(runErr, &exitErr) || exitErr.ExitCode() != 1) {
				t.Fatalf("depguard failed: %v, %s", runErr, stderr.String())
			}

			var report struct {
				Issues []struct {
					FromLinter string
					Pos        struct{ Filename string }
				}
			}
			if err := json.NewDecoder(bytes.NewReader(output)).Decode(&report); err != nil {
				t.Fatalf("decode linter output: %v, %s", err, output)
			}
			actual := make(map[string]bool)
			for _, issue := range report.Issues {
				if issue.FromLinter != "depguard" {
					t.Fatalf("unexpected linter %q", issue.FromLinter)
				}
				actual[issue.Pos.Filename] = true
			}
			for file, rejected := range expected {
				if actual[file] != rejected {
					t.Errorf("depguard %s: rejected = %t, expected %t", file, actual[file], rejected)
				}
				delete(actual, file)
			}
			if len(actual) != 0 {
				t.Errorf("unexpected depguard files: %v", actual)
			}
		})
	}
}

func writeGuardFixture(t *testing.T, root, name, contents string) {
	t.Helper()

	file := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(file), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(contents), 0o600); err != nil {
		t.Fatal(fmt.Errorf("write guard fixture: %w", err))
	}
}
