package parser_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/abad/ast"
	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/parser"
	"github.com/NeowayLabs/abad/token"
	"github.com/madlambda/spells/assert"
)

var E = fmt.Errorf

func TestParserNumbers(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "SmallDecimal",
			code: "1",
			want: ast.NewIntNumber(1),
		},
		{
			name: "BigDecimal",
			code: "1234567890",
			want: ast.NewIntNumber(1234567890),
		},
		{
			name:    "InvalidDecimal",
			code:    "1a",
			wantErr: E("tests.js:1:0: invalid token: 1a"),
		},
		{
			name: "SmallHexadecimal",
			code: "0x0",
			want: ast.NewIntNumber(0),
		},
		{
			name: "BigHexaDecimal",
			code: "0x1234567890abcdef",
			want: ast.NewIntNumber(0x1234567890abcdef),
		},
		{
			name: "HexadecimalFF",
			code: "0xff",
			want: ast.NewIntNumber(0xff),
		},
		{
			name: "SmallRealNumber",
			code: ".1",
			want: ast.NewNumber(0.1),
		},
		{
			name: "ZeroRealNumer",
			code: ".0000",
			want: ast.NewNumber(0.0),
		},
		{
			name: "SomeDecimal",
			code: "1234",
			want: ast.NewIntNumber(1234),
		},
		{
			name: "SmallRealNumberWithMultipleDigits",
			code: "0.12345",
			want: ast.NewNumber(0.12345),
		},
		{
			name:    "InvalidRealNumberWithLetter",
			code:    "0.a",
			wantErr: E("tests.js:1:0: invalid token: 0.a"),
		},
		{
			name:    "InvalidRealNumberWithTwoDots",
			code:    "12.13.",
			wantErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			name: "RealNumberWithExponent",
			code: "1.0e10",
			want: ast.NewNumber(1.0e10),
		},
		{
			name: "DecimalWithExponent",
			code: "1e10",
			want: ast.NewNumber(1e10),
		},
		{
			name: "SmallRealNumberWithExponent",
			code: ".1e10",
			want: ast.NewNumber(.1e10),
		},
		{
			name: "DecimalWithNegativeExponent",
			code: "1e-10",
			want: ast.NewNumber(1e-10),
		},
		{
			name: "NegativeDecimalWithOneDigit",
			code: "-1",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1),
			),
		},
		{
			name: "NegativeDecimalWithMultipleDigits",
			code: "-1234",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1234),
			),
		},
		{
			name: "NegativeZeroHexadecimal",
			code: "-0x0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name: "NegativeFFHexadecimal",
			code: "-0xff",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(255),
			),
		},
		{
			name: "NegativeZeroRealNumber",
			code: "-.0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name: "NegativeZeroRealNumberWithExponent",
			code: "-.0e1",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(0),
			),
		},
		{
			name:    "InvalidNegativeRealNumber",
			code:    "-12.13.",
			wantErr: E("tests.js:1:0: invalid token: 12.13."),
		},
		{
			name: "NegativeDecimalWithNegativeExponent",
			code: "-1e-10",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewNumber(1.0e-10),
			),
		},
		{
			name: "NegativePlusZeroDecimal",
			code: "-+0",
			want: ast.NewUnaryExpr(
				token.Minus, ast.NewUnaryExpr(
					token.Plus, ast.NewNumber(0),
				),
			),
		},
		{
			name: "InterleavedNegativeWithPlusAndZeroDecimal",
			code: "-+-+0",
			want: ast.NewUnaryExpr(token.Minus,
				ast.NewUnaryExpr(token.Plus,
					ast.NewUnaryExpr(token.Minus,
						ast.NewUnaryExpr(token.Plus,
							ast.NewNumber(0))))),
		},
	})
}

func TestIdentifier(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Underscore",
			code: "_",
			want: identifier("_"),
		},
		{
			name: "Dolar",
			code: "$",
			want: identifier("$"),
		},
		{
			name: "Console",
			code: "console",
			want: identifier("console"),
		},
		{
			name: "AngularSux",
			code: "angular",
			want: identifier("angular"),
		},
		{
			name: "HyperdUnderscores",
			code: "___hyped___",
			want: identifier("___hyped___"),
		},
		{
			name: "LettersAndDolars",
			code: "a$b$c",
			want: identifier("a$b$c"),
		},
		{
			name: "WithSemicolon",
			code: "a;",
			want: identifier("a"),
		},
		{
			name: "SeparatedBySemicolon",
			code: "a;b;c",
			wants: []ast.Node{
				identifier("a"),
				identifier("b"),
				identifier("c"),
			},
		},
	})
}

func TestString(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Empty",
			code: `""`,
			want: str(""),
		},
		{
			name: "JustSpace",
			code: `" "`,
			want: str(" "),
		},
		{
			name: "CommonName",
			code: `"inferno"`,
			want: str("inferno"),
		},
		{
			name: "LotsOfChars",
			code: `"!@#$%&*()]}[{/?^~ç"`,
			want: str("!@#$%&*()]}[{/?^~ç"),
		},
	})
}

func TestKeywords(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Null",
			code: "null",
			want: null(),
		},
		{
			name: "Undefined",
			code: "undefined",
			want: undefined(),
		},
		{
			name: "FalseBool",
			code: "false",
			want: boolean(false),
		},
		{
			name: "TrueBool",
			code: "true",
			want: boolean(true),
		},
	})
}

func TestMemberExpr(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "AccessingLogOnConsole",
			code: "console.log",
			want: memberExpr(identifier("console"), "log"),
		},
		{
			name:    "ErrorAccessingEmptyMember",
			code:    "console.",
			wantErr: E("tests.js:1:0: unexpected EOF"),
		},
		{
			name: "AccessMemberOfSelf",
			code: "self.a",
			want: memberExpr(identifier("self"), "a"),
		},
		{
			name: "OneLevelOfNesting",
			code: "self.self.self", // same as: (self.self).self
			want: memberExpr(
				memberExpr(identifier("self"), "self"),
				"self",
			),
		},
		{
			name: "MultipleLevelsOfNesting",
			code: "a.b.c.d.e.f", // same as: ((((a.b).c).d).e).f)
			want: memberExpr(
				memberExpr(
					memberExpr(
						memberExpr(
							memberExpr(identifier("a"), "b"),
							"c",
						),
						"d",
					),
					"e",
				),
				"f",
			),
		},
	})
}

func TestVarDeclarationErrors(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "NotIdentifierAfterVarDecl",
			code: "var 1234;",
			fail: true,
		},
		{
			name: "VarWithoutIdentifier",
			code: "var",
			fail: true,
		},
		{
			name: "EOFAfterIdentifier",
			code: "var x",
			fail: true,
		},
		{
			name: "EOFAfterInitializer",
			code: "var x =",
			fail: true,
		},
		{
			name: "InvalidAssignExpression",
			code: "var a = var;",
			fail: true,
		},
		{
			name: "InvalidAssignInitializer",
			code: "var a ! 5;",
			fail: true,
		},
		{
			name: "InvalidVarOnMultipleInits",
			code: "var a = 5, var x = 6;",
			fail: true,
		},
		{
			name: "InvalidAssignOnMultipleInits",
			code: "var a = 5, x ! 3",
			fail: true,
		},
		{
			name: "InvalidSeparatorForMultipleInits",
			code: "var d = 6 : x = 1",
			fail: true,
		},
		{
			name: "InvalidFuncall",
			code: "var d = lala(666",
			fail: true,
		},
		{
			name: "InvalidMemberAccess",
			code: "var d = obj.666",
			fail: true,
		},
	})
}

func TestVarStatement(t *testing.T) {
	// http://es5.github.io/#x12.2

	// WHY: just to avoid typing a lot initializing tests
	// that have single vars being initialized.
	vars := func(name ast.Ident, val ast.Node) ast.Node {
		return varDecls(varDecl(name, val))
	}

	// TODO: add vars init to funcall and access member expressions
	// eg: var a = func()
	//     var b = a.x.i()

	// TODO: add identifier to multiple vars statements and multiple vars on
	// single statement.
	runTests(t, []TestCase{
		{
			name: "NoInitializer",
			code: "var x;",
			want: vars(identifier("x"), undefined()),
		},
		{
			name: "Decimal",
			code: "var y = 1;",
			want: vars(identifier("y"), intNumber(1)),
		},
		{
			name: "Real",
			code: "var y = 6.66;",
			want: vars(identifier("y"), number(6.66)),
		},
		{
			name: "Hex",
			code: "var y = 0xFF;",
			want: vars(identifier("y"), intNumber(255)),
		},
		{
			name: "String",
			code: `var win = "i4k likes windows";`,
			want: vars(identifier("win"), str("i4k likes windows")),
		},
		{
			name: "Undefined",
			code: "var u = undefined;",
			want: vars(identifier("u"), undefined()),
		},
		{
			name: "Null",
			code: "var u = null;",
			want: vars(identifier("u"), null()),
		},
		{
			name: "True",
			code: "var b = true;",
			want: vars(identifier("b"), boolean(true)),
		},
		{
			name: "False",
			code: "var b = false;",
			want: vars(identifier("b"), boolean(false)),
		},
		{
			name: "Identifier",
			code: "var b = a;",
			want: vars(identifier("b"), identifier("a")),
		},
		{
			name: "MultipleVarsInSingleStatement",
			code: `
					var d = 666,
					    x = 0xFF,
					    s = "hi",
					    u = undefined,
					    n = null;
			`,
			want: varDecls(
				varDecl(identifier("d"), intNumber(666)),
				varDecl(identifier("x"), intNumber(255)),
				varDecl(identifier("s"), str("hi")),
				varDecl(identifier("u"), undefined()),
				varDecl(identifier("n"), null()),
			),
		},
		{
			name: "MultipleVarsStatements",
			code: `
					var d = 666;
					var x = 0xFF;
					var s = "hi";
					var u = undefined;
					var n = null;
			`,
			wants: []ast.Node{
				vars(identifier("d"), intNumber(666)),
				vars(identifier("x"), intNumber(255)),
				vars(identifier("s"), str("hi")),
				vars(identifier("u"), undefined()),
				vars(identifier("n"), null()),
			},
		},
	})
}

func TestParserFuncallError(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "NoParamMissingRParen",
			code: "a(",
			fail: true,
		},
		{
			name: "OneParamMissingRParen",
			code: "a(666",
			fail: true,
		},
		{
			name: "MultipleParamsMissingRParen",
			code: "a(666,777",
			fail: true,
		},
	})
}

func TestParserFuncall(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "NoParameter",
			code: "a()",
			want: callExpr(identifier("a"), []ast.Node{}),
		},
		{
			name: "UndefinedParameter",
			code: "b(undefined)",
			want: callExpr(identifier("b"), []ast.Node{undefined()}),
		},
		{
			name: "NullParameter",
			code: "b(null)",
			want: callExpr(identifier("b"), []ast.Node{null()}),
		},
		{
			name: "TrueBoolParameter",
			code: "b(true)",
			want: callExpr(identifier("b"), []ast.Node{boolean(true)}),
		},
		{
			name: "FalseBoolParameter",
			code: "b(false)",
			want: callExpr(identifier("b"), []ast.Node{boolean(false)}),
		},
		{
			name: "IntParameter",
			code: "b(1)",
			want: callExpr(identifier("b"), []ast.Node{intNumber(1)}),
		},
		{
			name: "HexParameter",
			code: "d(0xFF)",
			want: callExpr(identifier("d"), []ast.Node{intNumber(255)}),
		},
		{
			name: "NumberParameter",
			code: "c(6.66)",
			want: callExpr(identifier("c"), []ast.Node{number(6.66)}),
		},
		{
			name: "StringParameter",
			code: `c("hi")`,
			want: callExpr(identifier("c"), []ast.Node{str("hi")}),
		},
		{
			name: "MemberAccessWithoutParams",
			code: "console.log()",
			want: callExpr(
				memberExpr(identifier("console"), "log"),
				[]ast.Node{},
			),
		},
		{
			name: "MultipleCallsSplitByLotsOfSemiColons",
			code: "a();;;;;b();;",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{}),
				callExpr(identifier("b"), []ast.Node{}),
			},
		},
		{
			name: "MultipleCallsSplitByLotsOfSemiColonsNewlines",
			code: "a();\n\n;\n;;;b();\n;",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{}),
				callExpr(identifier("b"), []ast.Node{}),
			},
		},
		{
			name: "MultipleCallsSplitBySemiColon",
			code: "a();b();",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{}),
				callExpr(identifier("b"), []ast.Node{}),
			},
		},
		{
			name: "MultipleCallsSplitBySemiColonNewline",
			code: "a();\nb();",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{}),
				callExpr(identifier("b"), []ast.Node{}),
			},
		},
		{
			name: "MultipleCallsSplitBySemiColonWithParams",
			code: "a(1.1);b(0xFF);c(666);",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{number(1.1)}),
				callExpr(identifier("b"), []ast.Node{intNumber(255)}),
				callExpr(identifier("c"), []ast.Node{intNumber(666)}),
			},
		},
		{
			name: "MultipleCallsSplitBySemiColonNewlinesWithParams",
			code: "a(1.1);\nb(0xFF);\nc(666);",
			wants: []ast.Node{
				callExpr(identifier("a"), []ast.Node{number(1.1)}),
				callExpr(identifier("b"), []ast.Node{intNumber(255)}),
				callExpr(identifier("c"), []ast.Node{intNumber(666)}),
			},
		},
		{
			name: "MultipleMemberAccessSplitBySemicolon",
			code: "console.log(2.0);console.log(666);",
			wants: []ast.Node{
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{number(2.0)},
				),
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{intNumber(666)},
				),
			},
		},
		{
			name: "MultipleMemberAccessSplitBySemicolonNewline",
			code: "console.log(2.0);\nconsole.log(666);",
			wants: []ast.Node{
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{number(2.0)},
				),
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{intNumber(666)},
				),
			},
		},
		{
			name: "MemberAccessWithDecimalParam",
			code: "console.log(2.0)",
			want: callExpr(
				memberExpr(identifier("console"), "log"),
				[]ast.Node{ast.NewNumber(2.0)},
			),
		},
		{
			name: "NestedMemberAccessWithDecimalParam",
			code: "self.console.log(2.0)",
			want: callExpr(
				memberExpr(
					memberExpr(
						identifier("self"),
						"console",
					),
					"log",
				),
				[]ast.Node{ast.NewNumber(2.0)},
			),
		},
		{
			name: "AllTypesTogether",
			code: `all("hi",true,false,null,undefined,666,0xFF)`,
			want: callExpr(identifier("all"), []ast.Node{
				str("hi"),
				boolean(true),
				boolean(false),
				null(),
				undefined(),
				intNumber(666),
				intNumber(255),
			}),
		},
		{
			name: "AllTypesTogetherWithSpaces",
			code: `all( "hi", true, false, null, undefined, 666, 0xFF )`,
			want: callExpr(identifier("all"), []ast.Node{
				str("hi"),
				boolean(true),
				boolean(false),
				null(),
				undefined(),
				intNumber(666),
				intNumber(255),
			}),
		},
	})
}

func TestFunDecl(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "simple function",
			code: `function a(){}`,
			want: fundecl(
				identifier("a"),
				[]ast.Ident{},
				program(),
			),
		},
		{
			name: "function with args",
			code: `function a(b, c, d){}`,
			want: fundecl(
				identifier("a"),
				[]ast.Ident{identifier("b"), identifier("c"), identifier("d")},
				program(),
			),
		},
		{
			name: "function with args and body",
			code: `function a(b){b(1, 2)}`,
			want: fundecl(
				identifier("a"),
				[]ast.Ident{identifier("b")},
				program(
					callExpr(identifier("b"), []ast.Node{
						number(1), number(2),
					})),
			),
		},
		{
			name: "function between stmts",
			code: `console.log(1);
			function a(b){
				b(1, 2)
			}
			console.log(2);
			`,
			wants: []ast.Node{
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{ast.NewNumber(1.0)},
				),
				fundecl(
					identifier("a"),
					[]ast.Ident{identifier("b")},
					program(
						callExpr(identifier("b"), []ast.Node{
							number(1), number(2),
						})),
				),
				callExpr(
					memberExpr(identifier("console"), "log"),
					[]ast.Node{ast.NewNumber(2.0)},
				),
			},
		},
	})
}

// TestCase is the description of an parser related test.
// The fields want and wants are mutually exclusive, you should
// never provide both. If "wants" is provided the "want" field will be ignored.
//
// This is supposed to make it easier to test single nodes and multiple nodes.
type TestCase struct {
	name    string
	code    string
	want    ast.Node
	wants   []ast.Node
	fail    bool
	wantErr error
}

func (tc *TestCase) run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {
		tree, err := parser.Parse("tests.js", tc.code)

		if tc.fail && tc.wantErr == nil {
			if err == nil {
				t.Fatalf("expected an error got succeess:\n%s\n", tree)
			}
		} else {
			assert.EqualErrs(t, tc.wantErr, err, "parser err")
		}

		if err != nil {
			return
		}

		if tc.wants == nil {
			assertEqualNodes(t, []ast.Node{tc.want}, tree.Nodes)
			return
		}

		assertEqualNodes(t, tc.wants, tree.Nodes)
	})
}

func runTests(t *testing.T, tcases []TestCase) {
	for _, tcase := range tcases {
		tcase.run(t)
	}
}

func assertEqualNodes(t *testing.T, want []ast.Node, got []ast.Node) {
	if len(want) != len(got) {
		t.Errorf("want[%d] nodes but got[%d] nodes", len(want), len(got))
		t.Fatalf("want:\n%v\n\ngot:\n%v\n", want, got)
	}

	for i, w := range want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("wanted node[%d][%v] != got node[%d][%v]", i, w, i, g)
		}
	}
}

func number(n float64) ast.Number {
	return ast.NewNumber(n)
}

func intNumber(n int64) ast.Number {
	return ast.NewIntNumber(n)
}

func identifier(val string) ast.Ident {
	return ast.NewIdent(utf16.S(val))
}

func str(val string) ast.String {
	return ast.NewString(utf16.S(val))
}

func null() ast.Null {
	return ast.NewNull()
}

func undefined() ast.Undefined {
	return ast.NewUndefined()
}

func boolean(b bool) ast.Bool {
	return ast.NewBool(b)
}

func memberExpr(obj ast.Node, memberName string) *ast.MemberExpr {
	return ast.NewMemberExpr(obj, identifier(memberName))
}

func callExpr(callee ast.Node, args []ast.Node) *ast.CallExpr {
	return ast.NewCallExpr(callee, args)
}

func fundecl(name ast.Ident, args []ast.Ident, body *ast.Program) *ast.FunDecl {
	return ast.NewFunDecl(name, args, body)
}

func program(stmts ...ast.Node) *ast.Program {
	return &ast.Program{
		Nodes: stmts,
	}
}

func varDecls(vars ...ast.VarDecl) ast.VarDecls {
	return ast.NewVarDecls(vars...)
}

func varDecl(name ast.Ident, value ast.Node) ast.VarDecl {
	return ast.NewVarDecl(name, value)
}
