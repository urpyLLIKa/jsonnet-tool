package natives

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/go-jsonnet"

	"github.com/google/go-jsonnet/ast"
)

var errRegularExpression = errors.New("regular expression")

// escapeStringRegex escapes all regular expression metacharacters
// and returns a regular expression that matches the literal text.
func escapeStringRegex() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: ast.Identifiers{"str"},
		Func: func(s []interface{}) (interface{}, error) {
			str, ok := s[0].(string)
			if !ok {
				return "", fmt.Errorf("str must be a string: %w", errUnexpectedArgumentType)
			}

			result := regexp.QuoteMeta(str)

			return result, nil
		},
	}
}

// regexMatch returns whether the given string is matched by the given re2 regular expression.
func regexMatch() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: ast.Identifiers{"regex", "str"},
		Func: func(s []interface{}) (interface{}, error) {
			regex, ok := s[0].(string)
			if !ok {
				return false, fmt.Errorf("regex must be a string: %w", errUnexpectedArgumentType)
			}

			str, ok := s[1].(string)
			if !ok {
				return false, fmt.Errorf("str must be a string: %w", errUnexpectedArgumentType)
			}

			result, err := regexp.MatchString(regex, str)
			if err != nil {
				return false, errors.Join(err, errRegularExpression)
			}

			return result, nil
		},
	}
}

// regexSubst replaces all matches of the re2 regular expression with another string.
func regexSubst() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexSubst",
		Params: ast.Identifiers{"regex", "src", "repl"},
		Func: func(s []interface{}) (interface{}, error) {
			regex, ok := s[0].(string)
			if !ok {
				return "", fmt.Errorf("regex must be a string: %w", errUnexpectedArgumentType)
			}

			src, ok := s[1].(string)
			if !ok {
				return "", fmt.Errorf("str must be a string: %w", errUnexpectedArgumentType)
			}

			repl, ok := s[2].(string)
			if !ok {
				return "", fmt.Errorf("repl must be a string: %w", errUnexpectedArgumentType)
			}

			r, err := regexp.Compile(regex)
			if err != nil {
				return "", errors.Join(err, errRegularExpression)
			}

			return r.ReplaceAllString(src, repl), nil
		},
	}
}
