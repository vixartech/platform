// Command goarch validates Vixar hexagonal package boundaries.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

var errInvalidLocation = errors.New("invalid application package")

func run(ctx context.Context) error {
	applicationRoot := filepath.Join("internal", "app")
	if _, err := os.Stat(applicationRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("inspect application root: %w", err)
	}

	modulePath, err := commandOutput(ctx, "go", "list", "-m", "-f", "{{.Path}}")
	if err != nil {
		return err
	}

	packages, err := listPackages(ctx)
	if err != nil {
		return err
	}

	violations := validatePackages(strings.TrimSpace(modulePath), packages)
	if len(violations) != 0 {
		return fmt.Errorf("hexagonal architecture violations:\n%s", strings.Join(violations, "\n"))
	}

	return nil
}

func commandOutput(ctx context.Context, name string, arguments ...string) (string, error) {
	command := exec.CommandContext(ctx, name, arguments...)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("run %s: %w: %s", name, err, strings.TrimSpace(stderr.String()))
	}

	return string(output), nil
}

type packageInfo struct {
	ImportPath string
	Imports    []string
}

func listPackages(ctx context.Context) ([]packageInfo, error) {
	command := exec.CommandContext(ctx, "go", "list", "-json", "./internal/app/...")

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	command.Stderr = &stderr
	command.Stdout = &stdout
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("list application packages: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	decoder := json.NewDecoder(&stdout)
	packages := make([]packageInfo, 0)

	for {
		var current packageInfo
		if err := decoder.Decode(&current); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decode application package: %w", err)
		}

		packages = append(packages, current)
	}

	return packages, nil
}

func validatePackages(modulePath string, packages []packageInfo) []string {
	violations := make([]string, 0)

	for _, current := range packages {
		location, local, err := locate(modulePath, current.ImportPath)
		if err != nil {
			violations = append(violations, err.Error())

			continue
		}
		if !local {
			continue
		}

		for _, imported := range current.Imports {
			target, targetLocal, targetErr := locate(modulePath, imported)
			if targetErr != nil {
				violations = append(violations, targetErr.Error())

				continue
			}
			if !targetLocal || permits(location, target) {
				continue
			}

			violations = append(
				violations,
				fmt.Sprintf("- %s (%s) must not import %s (%s)", current.ImportPath, location.layer, imported, target.layer),
			)
		}
	}

	sort.Strings(violations)

	return violations
}

type packageLocation struct {
	service string
	layer   string
	owner   string
}

func locate(modulePath string, importPath string) (packageLocation, bool, error) {
	prefix := modulePath + "/internal/app/"
	if !strings.HasPrefix(importPath, prefix) {
		return packageLocation{}, false, nil
	}

	parts := strings.Split(strings.TrimPrefix(importPath, prefix), "/")
	if len(parts) < 2 {
		return packageLocation{}, true, fmt.Errorf(
			"%w: %s must be placed in domain, service, ports, or adapters",
			errInvalidLocation,
			importPath,
		)
	}

	location := packageLocation{service: parts[0], layer: parts[1]}
	if !knownLayer(location.layer) {
		return packageLocation{}, true, fmt.Errorf(
			"%w: %s uses unknown application layer %q",
			errInvalidLocation,
			importPath,
			location.layer,
		)
	}

	if (location.layer == "ports" || location.layer == "adapters") && len(parts) < 3 {
		return packageLocation{}, true, fmt.Errorf(
			"%w: %s must declare a concrete %s owner",
			errInvalidLocation,
			importPath,
			location.layer,
		)
	}
	if len(parts) >= 3 {
		location.owner = parts[2]
	}

	return location, true, nil
}

func knownLayer(layer string) bool {
	switch layer {
	case "domain", "service", "ports", "adapters":
		return true
	default:
		return false
	}
}

func permits(source packageLocation, target packageLocation) bool {
	if source.service != target.service {
		return false
	}

	switch source.layer {
	case "domain":
		return target.layer == "domain"
	case "service":
		return target.layer == "domain" || target.layer == "service"
	case "ports":
		return target.layer == "domain" || target.layer == "service" || sameBoundary(source, target)
	case "adapters":
		return target.layer == "domain" || target.layer == "service" || sameBoundary(source, target)
	default:
		return false
	}
}

func sameBoundary(source packageLocation, target packageLocation) bool {
	return source.layer == target.layer && source.owner == target.owner
}
