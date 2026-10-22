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
	"strconv"
	"strings"
	"testing"
)

type guardApp struct {
	app   string
	ports bool
}

type guardReport struct {
	Issues []guardIssue `json:"Issues"` //nolint:tagliatelle // External golangci-lint JSON field names.
}

type guardIssue struct {
	FromLinter string        `json:"FromLinter"` //nolint:tagliatelle // External golangci-lint JSON field names.
	Pos        guardPosition `json:"Pos"`        //nolint:tagliatelle // External golangci-lint JSON field names.
}

type guardPosition struct {
	Filename string `json:"Filename"` //nolint:tagliatelle // External golangci-lint JSON field names.
}

// The real linter verifies shared YAML matching, including arbitrary checkout names.
//
//nolint:paralleltest // Run seven real linter processes sequentially to bound CPU and memory use.
func TestQueryDepguard(t *testing.T) {
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

	for _, registered := range registeredGuardApps() {
		t.Run(registered.app, func(t *testing.T) {
			root := guardRoot(t, registered.app)
			expected := guardQueryFixtures(t, root, registered)
			addDirectGuards(t, root, registered, expected)
			addServiceGuards(t, root, registered.app, expected)
			addExternalGuards(t, root, registered.app, expected)
			vendorGuardFixture(t, binary, root, env)
			output := runGuardLinter(t, linter, config, root, env)
			assertGuardOutput(t, expected, output)
		})
	}
}

func registeredGuardApps() []guardApp {
	return []guardApp{
		{app: "gost", ports: true},
		{app: "dozor", ports: true},
		{app: "id", ports: true},
		{app: "mironenko", ports: true},
		{app: "signal", ports: true},
		{app: "relay", ports: false},
		{app: "srez", ports: true},
	}
}

func guardRoot(t *testing.T, app string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "arbitrary-checkout")
	module := "github.com/vixartech/" + app
	writeGuardFixture(t, root, "go.mod", "module "+module+"\n\ngo 1.26.5\n"+
		"\nrequire golang.org/x/net v0.0.0\nreplace golang.org/x/net => ./xnet\n")

	actual, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}

	return actual
}

func guardQueryFixtures(t *testing.T, root string, registered guardApp) map[string]bool {
	t.Helper()

	module := "github.com/vixartech/" + registered.app
	expected := make(map[string]bool)

	for index, test := range queryImportCases() {
		if strings.HasPrefix(test.source, "ports/") && !registered.ports {
			continue
		}

		source := "internal/app/" + registered.app + "/" + test.source + "/case" + strconv.Itoa(index) + "/check.go"
		target := strings.ReplaceAll(test.target, testModule, module)
		target = strings.ReplaceAll(target, "/internal/app/example/", "/internal/app/"+registered.app+"/")
		writeGuardFixture(t, root, source, "package fixture\nimport _ "+strconv.Quote(target)+"\n")

		if local, ok := strings.CutPrefix(target, module+"/"); ok {
			writeGuardFixture(t, root, local+"/empty.go", "package fixture\n")
		}

		expected[filepath.Join(root, filepath.FromSlash(source))] = !test.allowed
	}

	return expected
}

func directGuardCases() []queryImportCase {
	return []queryImportCase{
		{source: "query", target: "ports/grpc"},
		{source: "domain", target: "query"},
		{source: "service", target: "query"},
		{source: "ports/http", target: "query", allowed: true},
		{source: "adapters/postgres", target: "query", allowed: true},
	}
}

func addDirectGuards(t *testing.T, root string, registered guardApp, expected map[string]bool) {
	t.Helper()

	module := "github.com/vixartech/" + registered.app
	// Distinct targets keep the direct-file fixture acyclic.
	for _, direct := range directGuardCases() {
		if strings.HasPrefix(direct.source, "ports/") && !registered.ports {
			continue
		}

		file := "internal/app/" + registered.app + "/" + direct.source + "/direct.go"
		target := module + "/internal/app/" + registered.app + "/" + direct.target
		writeGuardFixture(t, root, file, "package fixture\nimport _ "+strconv.Quote(target)+"\n")
		writeGuardFixture(t, root, strings.TrimPrefix(target, module+"/")+"/empty.go", "package fixture\n")
		expected[filepath.Join(root, filepath.FromSlash(file))] = !direct.allowed
	}
}

func addServiceGuards(t *testing.T, root, app string, expected map[string]bool) {
	t.Helper()

	fixture := "github.com/vixartech/" + app + "/internal/test/fixture"

	writeGuardFixture(t, root, "internal/test/fixture/empty.go", "package fixture\n")

	for _, file := range []string{"fixture.go", "fixture_test.go"} {
		path := "internal/app/" + app + "/service/" + file
		writeGuardFixture(t, root, path, "package fixture\nimport _ "+strconv.Quote(fixture)+"\n")
		expected[filepath.Join(root, filepath.FromSlash(path))] = file == "fixture.go" || app != "srez"
	}
}

func addExternalGuards(t *testing.T, root, app string, expected map[string]bool) {
	t.Helper()
	writeGuardFixture(t, root, "xnet/go.mod", "module golang.org/x/net\n\ngo 1.26.5\n")

	imports := []string{"golang.org/x/net/idna", "golang.org/x/net/idna/subpkg", "golang.org/x/net/http2"}
	for index, imported := range imports {
		file := "internal/app/" + app + "/domain/external" + strconv.Itoa(index) + ".go"
		writeGuardFixture(t, root, file, "package fixture\nimport _ "+strconv.Quote(imported)+"\n")
		local := strings.TrimPrefix(imported, "golang.org/x/net/")
		writeGuardFixture(t, root, "xnet/"+local+"/empty.go", "package fixture\n")
		expected[filepath.Join(root, filepath.FromSlash(file))] = imported != "golang.org/x/net/idna"
	}
}

func guardGoBinary(t *testing.T) string {
	t.Helper()
	command := exec.CommandContext(t.Context(), "go", "env", "GOROOT")

	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(strings.TrimSpace(string(output)), "bin", "go")
}

func guardEnvironment(binary string) []string {
	return append(os.Environ(), "GO111MODULE=on", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off",
		"PATH="+filepath.Dir(binary)+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func vendorGuardFixture(t *testing.T, binary, root string, env []string) {
	t.Helper()

	command := exec.CommandContext(t.Context(), binary, "mod", "vendor")
	command.Dir = root

	command.Env = env
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("vendor local guard fixture: %v, %s", err, output)
	}
}

func runGuardLinter(t *testing.T, linter, config, root string, env []string) []byte {
	t.Helper()
	command := exec.CommandContext(t.Context(), linter, "run", "--config", config,
		"--allow-parallel-runners", "--concurrency=2",
		"--enable-only", "depguard", "--output.json.path=stdout", "--output.text.path=", "./...")
	command.Dir = root
	command.Env = env

	var stderr bytes.Buffer

	command.Stderr = &stderr
	output, runErr := command.Output()

	var exitErr *exec.ExitError
	if runErr != nil && (!errors.As(runErr, &exitErr) || exitErr.ExitCode() != 1) {
		t.Fatalf("depguard failed: %v, %s", runErr, stderr.String())
	}

	return output
}

func assertGuardOutput(t *testing.T, expected map[string]bool, output []byte) {
	t.Helper()

	var report guardReport
	if err := json.NewDecoder(bytes.NewReader(output)).Decode(&report); err != nil {
		t.Fatalf("decode linter output: %v, %s", err, output)
	}

	actual := guardRejectedFiles(t, report)
	for file, rejected := range expected {
		if actual[file] != rejected {
			t.Errorf("depguard %s: rejected = %t, expected %t", file, actual[file], rejected)
		}

		delete(actual, file)
	}

	if len(actual) != 0 {
		t.Errorf("unexpected depguard files: %v", actual)
	}
}

func guardRejectedFiles(t *testing.T, report guardReport) map[string]bool {
	t.Helper()

	actual := make(map[string]bool)

	for _, issue := range report.Issues {
		if issue.FromLinter != "depguard" {
			t.Fatalf("unexpected linter %q", issue.FromLinter)
		}

		actual[issue.Pos.Filename] = true
	}

	return actual
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
