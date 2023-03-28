package checksum

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

// Parse recursively parses a jsonnet file and its imports, populating
// a map from filenames to sha256 hashes, to be used for checksumming.
func Parse(filename string, importer jsonnet.Importer, seen map[string][sha256.Size]byte) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %w", err)
	}

	root, err := jsonnet.SnippetToAST(filename, string(data))
	if err != nil {
		return fmt.Errorf("jsonnet.SnippetToAST: %w", err)
	}

	sum := sha256.Sum256(data)
	seen[filename] = sum

	var newFiles []string

	err = walk(root, func(node ast.Node) error {
		importNode, ok := node.(*ast.Import)
		if !ok {
			return nil
		}

		_, foundAt, importErr := importer.Import(filename, importNode.File.Value)
		if importErr != nil {
			return fmt.Errorf("importer.Import: %w", importErr)
		}

		if _, ok := seen[foundAt]; !ok {
			newFiles = append(newFiles, foundAt)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, newFile := range newFiles {
		err := Parse(newFile, importer, seen)
		if err != nil {
			return err
		}
	}

	return nil
}

func walk(node ast.Node, f func(ast.Node) error) error {
	err := f(node)
	if err != nil {
		return err
	}

	for _, child := range nodeChildren(node) {
		err = walk(child, f)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:funlen,gocognit,gocyclo,cyclop,goerr113,wsl,nlreturn
func nodeChildren(node ast.Node) []ast.Node {
	switch v := node.(type) {
	case *ast.Apply:
		var children []ast.Node
		children = append(children, v.Target)
		for _, named := range v.Arguments.Named {
			children = append(children, named.Arg)
		}
		for _, positional := range v.Arguments.Positional {
			children = append(children, positional.Expr)
		}
		return children
	case *ast.Array:
		var children []ast.Node
		for _, element := range v.Elements {
			children = append(children, element.Expr)
		}
		return children
	case *ast.ArrayComp:
		var children []ast.Node
		spec := &v.Spec
		for spec != nil {
			children = append(children, childrenFromForSpec(spec)...)
			spec = spec.Outer
		}
		return children
	case *ast.Assert:
		return []ast.Node{v.Cond, v.Message, v.Rest}
	case *ast.Binary:
		return []ast.Node{v.Left, v.Right}
	case *ast.Conditional:
		return []ast.Node{v.Cond, v.BranchTrue, v.BranchFalse}
	case *ast.DesugaredObject:
		var children []ast.Node
		children = append(children, v.Asserts...)
		for _, field := range v.Fields {
			children = append(children, field.Name, field.Body)
		}
		for _, local := range v.Locals {
			children = append(children, local.Body)
		}
		return children
	case *ast.Dollar:
		return []ast.Node{}
	case *ast.Error:
		return []ast.Node{v.Expr}
	case *ast.Function:
		var children []ast.Node
		children = append(children, v.Body)
		for _, parameter := range v.Parameters {
			children = append(children, parameter.DefaultArg)
		}
		return children
	case *ast.Import:
		return []ast.Node{}
	case *ast.ImportBin:
		return []ast.Node{}
	case *ast.ImportStr:
		return []ast.Node{}
	case *ast.InSuper:
		return []ast.Node{v.Index}
	case *ast.Index:
		return []ast.Node{v.Target, v.Index}
	case *ast.LiteralBoolean:
		return []ast.Node{}
	case *ast.LiteralNull:
		return []ast.Node{}
	case *ast.LiteralNumber:
		return []ast.Node{}
	case *ast.LiteralString:
		return []ast.Node{}
	case *ast.Local:
		var children []ast.Node
		children = append(children, v.Body)
		for _, bind := range v.Binds {
			children = append(children, bind.Body)
		}
		return children
	case *ast.Object:
		var children []ast.Node
		for _, field := range v.Fields {
			children = append(children, field.Expr1, field.Expr2, field.Expr3)
		}
		return children
	case *ast.ObjectComp:
		var children []ast.Node
		spec := &v.Spec
		for spec != nil {
			children = append(children, childrenFromForSpec(spec)...)
			spec = spec.Outer
		}
		for _, field := range v.Fields {
			children = append(children, field.Expr1, field.Expr2, field.Expr3)
		}
		return children
	case *ast.Parens:
		return []ast.Node{v.Inner}
	case *ast.Self:
		return []ast.Node{}
	case *ast.Slice:
		return []ast.Node{v.Target, v.BeginIndex, v.EndIndex, v.Step}
	case *ast.SuperIndex:
		return []ast.Node{v.Index}
	case *ast.Unary:
		return []ast.Node{v.Expr}
	case *ast.Var:
		return []ast.Node{}
	case nil:
		return []ast.Node{}
	default:
		panic(fmt.Errorf("nodeChildren does not support type %T", v))
	}
}

func childrenFromForSpec(v *ast.ForSpec) []ast.Node {
	children := make([]ast.Node, 0, len(v.Conditions)+1)

	children = append(children, v.Expr)
	for _, ifSpec := range v.Conditions {
		children = append(children, ifSpec.Expr)
	}

	return children
}
