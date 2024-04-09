package natives

import (
	"errors"

	"github.com/google/go-jsonnet"
)

var errUnexpectedArgumentType = errors.New("unexpected argument type")

var funcs = []*jsonnet.NativeFunction{
	// Regular expressions
	escapeStringRegex(),
	regexMatch(),
	regexSubst(),

	// Semver handling
	semverParse(),
	semverMatchesConstraint(),
}

// Register will register the native functions with the VM.
func Register(vm *jsonnet.VM) {
	for _, v := range funcs {
		vm.NativeFunction(v)
	}
}
