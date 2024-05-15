package natives

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

func semverParse() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "semverParse",
		Params: ast.Identifiers{"v"},
		Func: func(s []interface{}) (interface{}, error) {
			v, ok := s[0].(string)
			if !ok {
				return "", fmt.Errorf("v must be a string: %w", errUnexpectedArgumentType)
			}

			result, err := semver.NewVersion(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse semver: %w", err)
			}

			// uint64 is not a json type, convert to float64 for json
			parsedSemver := map[string]any{
				"major":      float64(result.Major()),
				"minor":      float64(result.Minor()),
				"patch":      float64(result.Patch()),
				"prerelease": result.Prerelease(),
				"metadata":   result.Metadata(),
			}

			return parsedSemver, nil
		},
	}
}

func semverMatchesConstraint() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "semverMatchesConstraint",
		Params: ast.Identifiers{"v", "constraint"},
		Func: func(s []interface{}) (interface{}, error) {
			v, ok := s[0].(string)
			if !ok {
				return "", fmt.Errorf("v must be a string: %w", errUnexpectedArgumentType)
			}

			constraint, ok := s[1].(string)
			if !ok {
				return "", fmt.Errorf("constraint must be a string: %w", errUnexpectedArgumentType)
			}

			sv, err := semver.NewVersion(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse semver: %w", err)
			}

			sc, err := semver.NewConstraint(constraint)
			if err != nil {
				return nil, fmt.Errorf("failed to parse constraint: %w", err)
			}

			matches, _ := sc.Validate(sv)

			return matches, nil
		},
	}
}
