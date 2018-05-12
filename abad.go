package abad

import (
	"fmt"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/token"
	"github.com/NeowayLabs/abad/types"
)

type (
	// Abad interpreter, a very bad one.
	Abad struct {
		filename string

		vars map[string]types.Value
	}
)

func NewAbad(filename string) *Abad {
	return &Abad{
		filename: filename,
		vars: map[string]types.Value{
			// TODO(i4k): remove
			// used to test reference until var decl is implemented
			"a": types.Number(666.),
		},
	}
}

func (a *Abad) Eval(code string) (types.Value, error) {
	program, err := parser.Parse(a.filename, code)
	if err != nil {
		return nil, err
	}

	return a.eval(program)
}

func (a *Abad) eval(n ast.Node) (types.Value, error) {
	switch n.Type() {
	case ast.NodeProgram:
		return a.evalProgram(n.(*ast.Program))
	case ast.NodeNumber:
		val := n.(ast.Number)
		return types.Number(val.Value()), nil
	case ast.NodeIdent:
		val := n.(ast.Ident)
		return a.evalIdentExpr(val)
	case ast.NodeUnaryExpr:
		expr := n.(*ast.UnaryExpr)
		return a.evalUnaryExpr(expr)
	}

	panic(fmt.Sprintf("AST(%s) not implemented", n))
	return nil, nil
}

func (a *Abad) evalProgram(stmts *ast.Program) (types.Value, error) {
	var (
		result types.Value
		err    error
	)
	for _, node := range stmts.Nodes {
		result, err = a.eval(node)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *Abad) evalUnaryExpr(expr *ast.UnaryExpr) (types.Value, error) {
	op := expr.Operator
	obj, err := a.eval(expr.Operand)
	if err != nil {
		return nil, err
	}

	// TODO(i4k): UnaryExpr could work in any expression in js
	// examples below are valid:
	//   -[]
	//	 +[]
	//   -{}
	//   +{}
	//   -eval("0")
	//   -new Object()
	num, ok := obj.(types.Number)
	if !ok {
		return nil, fmt.Errorf("not a number: %s", obj)
	}

	switch op {
	case token.Minus:
		num = -num
	case token.Plus:
		num = +num
	default:
		return nil, fmt.Errorf("unsupported unary operator: %s", op)
	}

	return num, nil
}

func (a *Abad) evalIdentExpr(ident ast.Ident) (types.Value, error) {
	val, ok := a.vars[ident.String()]
	if !ok {
		return nil, fmt.Errorf("%s is not defined", ident)
	}

	return val, nil
}