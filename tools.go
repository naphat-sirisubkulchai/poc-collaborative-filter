//go:build tools
// +build tools

// This file ensures that `go mod tidy` doesn't remove tool dependencies
package tools

import (
	_ "github.com/99designs/gqlgen" //nolint:typecheck
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
