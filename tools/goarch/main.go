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

const (
	layerDomain   = "domain"
	layerService  = "service"
	layerPorts    = "ports"
	layerAdapters = "adapters"

	boundaryOwnerIndex = 2
	boundaryParts      = boundaryOwnerIndex + 1
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

var (
	errArchitecture    = errors.New("hexagonal architecture violations")
	errInvalidLocation = errors.New("invalid application package")
	errMissingPath     = errors.New("application package has no import path")
)

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
		return fmt.Errorf("%w:\n%s", errArchitecture, strings.Join(violations, "\n"))
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
	output, err := commandOutput(ctx, "go", "list", "-json", "./internal/app/...")
	if err != nil {
		return nil, fmt.Errorf("list application packages: %w", err)
	}

	return decodePackages(strings.NewReader(output))
}

func decodePackages(input io.Reader) ([]packageInfo, error) {
	decoder := json.NewDecoder(input)
	packages := make([]packageInfo, 0)

	for {
		var fields map[string]json.RawMessage
		if err := decoder.Decode(&fields); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decode application package: %w", err)
		}

		current, err := decodePackage(fields)
		if err != nil {
			return nil, err
		}

		packages = append(packages, current)
	}

	return packages, nil
}

func decodePackage(fields map[string]json.RawMessage) (packageInfo, error) {
	pathData, exists := fields["ImportPath"]
	if !exists {
		return packageInfo{}, errMissingPath
	}

	var current packageInfo

	if err := json.Unmarshal(pathData, &current.ImportPath); err != nil {
		return packageInfo{}, fmt.Errorf("decode application import path: %w", err)
	}

	importsData, exists := fields["Imports"]
	if !exists {
		return current, nil
	}

	if err := json.Unmarshal(importsData, &current.Imports); err != nil {
		return packageInfo{}, fmt.Errorf("decode application imports: %w", err)
	}

	return current, nil
}

func validatePackages(modulePath string, packages []packageInfo) []string {
	violations := make([]string, 0, len(packages))

	for _, current := range packages {
		violations = append(violations, validatePackage(modulePath, current)...)
	}

	sort.Strings(violations)

	return violations
}

func validatePackage(modulePath string, current packageInfo) []string {
	location, local, err := locate(modulePath, current.ImportPath)
	if err != nil {
		return []string{err.Error()}
	}

	if !local {
		return nil
	}

	violations := make([]string, 0, len(current.Imports))

	for _, imported := range current.Imports {
		if violation := validateImport(modulePath, current.ImportPath, location, imported); violation != "" {
			violations = append(violations, violation)
		}
	}

	return violations
}

func validateImport(
	modulePath string,
	sourcePath string,
	source packageLocation,
	importPath string,
) string {
	target, local, err := locate(modulePath, importPath)
	if err != nil {
		return err.Error()
	}

	if !local || permits(source, target) {
		return ""
	}

	return fmt.Sprintf("- %s (%s) must not import %s (%s)", sourcePath, source.layer, importPath, target.layer)
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

	owner, err := boundaryOwner(location.layer, parts, importPath)
	if err != nil {
		return packageLocation{}, true, err
	}

	location.owner = owner

	return location, true, nil
}

func knownLayer(layer string) bool {
	switch layer {
	case layerDomain, layerService, layerPorts, layerAdapters:
		return true
	default:
		return false
	}
}

func boundaryOwner(layer string, parts []string, importPath string) (string, error) {
	if layer != layerPorts && layer != layerAdapters {
		return "", nil
	}

	if len(parts) < boundaryParts {
		return "", fmt.Errorf(
			"%w: %s must declare a concrete %s owner",
			errInvalidLocation,
			importPath,
			layer,
		)
	}

	return parts[boundaryOwnerIndex], nil
}

func permits(source packageLocation, target packageLocation) bool {
	if source.service != target.service {
		return false
	}

	switch source.layer {
	case layerDomain:
		return target.layer == layerDomain
	case layerService:
		return coreLayer(target.layer)
	case layerPorts, layerAdapters:
		return coreLayer(target.layer) || sameBoundary(source, target)
	default:
		return false
	}
}

func coreLayer(layer string) bool {
	switch layer {
	case layerDomain, layerService:
		return true
	default:
		return false
	}
}

func sameBoundary(source packageLocation, target packageLocation) bool {
	return source.layer == target.layer && source.owner == target.owner
}
