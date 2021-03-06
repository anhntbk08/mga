package handler

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

// Event describes an event struct.
type Event = gentypes.TypeRef

// Parse parses a given package, looks for a struct and returns it as a normalized structure.
func Parse(dir string, eventName string) (Event, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return Event{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(eventName)
		if obj == nil {
			continue
		}

		return ParseEvent(obj)
	}

	return Event{}, errors.New("event not found")
}

// ParseEvent parses an object as an event.
func ParseEvent(obj types.Object) (Event, error) {
	_, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return Event{}, fmt.Errorf("%q is not a struct", obj.Name())
	}

	event := Event{
		Name: obj.Name(),
		Package: gentypes.PackageRef{
			Name: obj.Pkg().Name(),
			Path: obj.Pkg().Path(),
		},
	}

	return event, nil
}
