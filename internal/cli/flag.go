/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package cli

import (
	"fmt"
	"reflect"
)

// Flag represents a command line flag configuration.
type Flag struct {
	Name        string
	Variable    string
	Type        string
	Default     any
	Description string
}

// typeKindMap maps string type names to their reflect.Kind counterparts.
var typeKindMap = map[string]reflect.Kind{
	"string": reflect.String,
	"int":    reflect.Int,
	"bool":   reflect.Bool,
}

// NewFlag creates a Flag instance with validation.
// It returns an error if the flag configuration is invalid.
func NewFlag(flag Flag) (*Flag, error) {
	var validationErrors []string
	if flag.Name == "" {
		validationErrors = append(validationErrors,
			fmt.Sprintf("Invalid name for `%s` flag", flag.Name))
	}
	if flag.Variable == "" {
		validationErrors = append(validationErrors,
			fmt.Sprintf("Invalid variable for `%s` flag", flag.Name))
	}
	if _, exists := typeKindMap[flag.Type]; !exists {
		validationErrors = append(validationErrors,
			fmt.Sprintf("Unsupported type for `%s` flag", flag.Name))
	}
	if reflect.TypeOf(flag.Default).Kind() != typeKindMap[flag.Type] {
		validationErrors = append(validationErrors,
			fmt.Sprintf("Invalid default value for `%s` flag", flag.Name))
	}
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("NewFlag(): %v", validationErrors)
	}
	return &flag, nil
}
