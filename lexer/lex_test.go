package lexer_test

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/NeowayLabs/abad/internal/utf16"
	"github.com/NeowayLabs/abad/lexer"
	"github.com/NeowayLabs/abad/token"
)

type TestCase struct {
	name          string
	code          utf16.Str
	want          []lexer.Tokval
	checkPosition bool
}

var Str func(string) utf16.Str = utf16.S
var EOF lexer.Tokval = lexer.EOF

func TestNumericLiterals(t *testing.T) {

	// SPEC: https://es5.github.io/#x7.8.3

	cases := []TestCase{
		{
			name: "SingleZero",
			code: Str("0"),
			want: tokens(decimalToken("0")),
		},
		{
			name: "BigDecimal",
			code: Str("1236547987794465977"),
			want: tokens(decimalToken("1236547987794465977")),
		},
		{
			name: "RealDecimalStartingWithPoint",
			code: Str(".1"),
			want: tokens(decimalToken(".1")),
		},
		{
			name: "RealDecimalEndingWithPoint",
			code: Str("1."),
			want: tokens(decimalToken("1.")),
		},
		{
			name: "LargeRealDecimalStartingWithPoint",
			code: Str(".123456789"),
			want: tokens(decimalToken(".123456789")),
		},
		{
			name: "SmallRealDecimal",
			code: Str("1.6"),
			want: tokens(decimalToken("1.6")),
		},
		{
			name: "BigRealDecimal",
			code: Str("11223243554.63445465789"),
			want: tokens(decimalToken("11223243554.63445465789")),
		},
		{
			name: "SmallRealDecimalWithSmallExponent",
			code: Str("1.0e1"),
			want: tokens(decimalToken("1.0e1")),
		},
		{
			name: "SmallDecimalWithSmallExponent",
			code: Str("1e1"),
			want: tokens(decimalToken("1e1")),
		},
		{
			name: "SmallDecimalWithSmallExponentUpperExponent",
			code: Str("1E1"),
			want: tokens(decimalToken("1E1")),
		},
		{
			name: "BigDecimalWithBigExponent",
			code: Str("666666666666e668"),
			want: tokens(decimalToken("666666666666e668")),
		},
		{
			name: "BigDecimalWithBigExponentUpperExponent",
			code: Str("666666666666E668"),
			want: tokens(decimalToken("666666666666E668")),
		},
		{
			name: "BigRealDecimalWithBigExponent",
			code: Str("666666666666.0e66"),
			want: tokens(decimalToken("666666666666.0e66")),
		},
		{
			name: "RealDecimalWithSmallNegativeExponent",
			code: Str("1.0e-1"),
			want: tokens(decimalToken("1.0e-1")),
		},
		{
			name: "RealDecimalWithBigNegativeExponent",
			code: Str("1.0e-50"),
			want: tokens(decimalToken("1.0e-50")),
		},
		{
			name: "SmallRealDecimalWithSmallUpperExponent",
			code: Str("1.0E1"),
			want: tokens(decimalToken("1.0E1")),
		},
		{
			name: "BigRealDecimalWithBigUpperExponent",
			code: Str("666666666666.0E66"),
			want: tokens(decimalToken("666666666666.0E66")),
		},
		{
			name: "RealDecimalWithSmallNegativeUpperExponent",
			code: Str("1.0E-1"),
			want: tokens(decimalToken("1.0E-1")),
		},
		{
			name: "RealDecimalWithBigNegativeUpperExponent",
			code: Str("1.0E-50"),
			want: tokens(decimalToken("1.0E-50")),
		},
		{
			name: "StartWithDotUpperExponent",
			code: Str(".0E-50"),
			want: tokens(decimalToken(".0E-50")),
		},
		{
			name: "StartWithDotExponent",
			code: Str(".0e5"),
			want: tokens(decimalToken(".0e5")),
		},
		{
			name: "ZeroHexadecimal",
			code: Str("0x0"),
			want: tokens(hexToken("0x0")),
		},
		{
			name: "BigHexadecimal",
			code: Str("0x123456789abcdef"),
			want: tokens(hexToken("0x123456789abcdef")),
		},
		{
			name: "BigHexadecimalUppercase",
			code: Str("0x123456789ABCDEF"),
			want: tokens(hexToken("0x123456789ABCDEF")),
		},
		{
			name: "LettersOnlyHexadecimal",
			code: Str("0xabcdef"),
			want: tokens(hexToken("0xabcdef")),
		},
		{
			name: "LettersOnlyHexadecimalUppercase",
			code: Str("0xABCDEF"),
			want: tokens(hexToken("0xABCDEF")),
		},
		{
			name: "ZeroHexadecimalUpperX",
			code: Str("0X0"),
			want: tokens(hexToken("0X0")),
		},
		{
			name: "BigHexadecimalUpperX",
			code: Str("0X123456789abcdef"),
			want: tokens(hexToken("0X123456789abcdef")),
		},
		{
			name: "BigHexadecimalUppercaseUpperX",
			code: Str("0X123456789ABCDEF"),
			want: tokens(hexToken("0X123456789ABCDEF")),
		},
		{
			name: "LettersOnlyHexadecimalUpperX",
			code: Str("0Xabcdef"),
			want: tokens(hexToken("0Xabcdef")),
		},
		{
			name: "LettersOnlyHexadecimalUppercaseUpperX",
			code: Str("0XABCDEF"),
			want: tokens(hexToken("0XABCDEF")),
		},
	}

	plusSignedCases := prependOnTestCases(TestCase{
		name: "PlusSign",
		code: Str("+"),
		want: []lexer.Tokval{plusToken()},
	}, cases)

	minusSignedCases := prependOnTestCases(TestCase{
		name: "MinusSign",
		code: Str("-"),
		want: []lexer.Tokval{minusToken()},
	}, cases)

	plusMinusPlusMinusSignedCases := prependOnTestCases(TestCase{
		name: "PlusMinusPlusMinusSign",
		code: Str("+-+-"),
		want: []lexer.Tokval{
			plusToken(),
			minusToken(),
			plusToken(),
			minusToken(),
		},
	}, cases)

	minusPlusMinusPlusSignedCases := prependOnTestCases(TestCase{
		name: "MinusPlusMinusPlusSign",
		code: Str("-+-+"),
		want: []lexer.Tokval{
			minusToken(),
			plusToken(),
			minusToken(),
			plusToken(),
		},
	}, cases)

	runTests(t, cases)
	runTests(t, plusSignedCases)
	runTests(t, minusSignedCases)
	runTests(t, plusMinusPlusMinusSignedCases)
	runTests(t, minusPlusMinusPlusSignedCases)
}

func TestStrings(t *testing.T) {
	// TODO: multiline strings
	// - escaped double quotes
	runTests(t, []TestCase{
		{
			name: "Empty",
			code: Str(`""`),
			want: tokens(stringToken("")),
		},
		{
			name: "SpacesOnly",
			code: Str(`"  "`),
			want: tokens(stringToken("  ")),
		},
		{
			name: "SingleChar",
			code: Str(`"k"`),
			want: tokens(stringToken("k")),
		},
		{
			name: "LotsOfCrap",
			code: Str(`"1234567890-+=abcdefg${[]})(()%_ /|/ yay %xi4klindaum"`),
			want: tokens(stringToken("1234567890-+=abcdefg${[]})(()%_ /|/ yay %xi4klindaum")),
		},
	})
}

func TestLineTerminator(t *testing.T) {
	type LineTerminator struct {
		name string
		val  string
	}

	lineTerminators := []LineTerminator{
		{name: "LineFeed", val: "\u000A"},
		{name: "CarriageReturn", val: "\u000D"},
		{name: "LineSeparator", val: "\u2028"},
		{name: "ParagraphSeparator", val: "\u2029"},
	}

	for _, lineTerminator := range lineTerminators {
		t.Run(lineTerminator.name, func(t *testing.T) {
			lt := lineTerminator.val
			runTests(t, []TestCase{
				{
					name: "Strings",
					code: sfmt(`"first"%s"second"`, lt),
					want: tokens(stringToken("first"), ltToken(lt), stringToken("second")),
				},
				{
					name: "Decimals",
					code: sfmt("1%s2", lt),
					want: tokens(decimalToken("1"), ltToken(lt), decimalToken("2")),
				},
				{
					name: "Hexadecimals",
					code: sfmt("0xFF%s0x11", lt),
					want: tokens(hexToken("0xFF"), ltToken(lt), hexToken("0x11")),
				},
				{
					name: "Identifiers",
					code: sfmt("hi%shello", lt),
					want: tokens(identToken("hi"), ltToken(lt), identToken("hello")),
				},
			})
		})
	}
}

func TestInvalidStrings(t *testing.T) {
	// TODO: add newline tests

	runTests(t, []TestCase{
		{
			name: "SingleDoubleQuote",
			code: Str(`"`),
			want: []lexer.Tokval{illegalToken(`"`)},
		},
		{
			name: "NoEndingDoubleQuote",
			code: Str(`"dsadasdsa123456`),
			want: []lexer.Tokval{illegalToken(`"dsadasdsa123456`)},
		},
	})
}

func TestIdentifiers(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "Underscore",
			code: Str("_"),
			want: tokens(identToken("_")),
		},
		{
			name: "SingleLetter",
			code: Str("a"),
			want: tokens(identToken("a")),
		},
		{
			name: "Self",
			code: Str("self"),
			want: tokens(identToken("self")),
		},
		{
			name: "Console",
			code: Str("console"),
			want: tokens(identToken("console")),
		},
		{
			name: "LotsUnderscores",
			code: Str("___hyped___"),
			want: tokens(identToken("___hyped___")),
		},
		{
			name: "DollarsInterwined",
			code: Str("a$b$c"),
			want: tokens(identToken("a$b$c")),
		},
		{
			name: "NumbersInterwined",
			code: Str("a1b2c"),
			want: tokens(identToken("a1b2c")),
		},
		{
			name: "AccessingMember",
			code: Str("console.log"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
			),
		},
		{
			name: "AccessingNoMember",
			code: Str("console."),
			want: tokens(
				identToken("console"),
				dotToken(),
			),
		},
		{
			name: "AccessingMemberOfMember",
			code: Str("console.log.toString"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
				dotToken(),
				identToken("toString"),
			),
		},
	})
}

func TestFuncall(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "OneLetterFunction",
			code: Str("a()"),
			want: tokens(
				identToken("a"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "BigFunctionName",
			code: Str("veryBigFunctionNameThatWouldAnnoyNatel()"),
			want: tokens(
				identToken("veryBigFunctionNameThatWouldAnnoyNatel"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "MemberFunction",
			code: Str("console.log()"),
			want: tokens(
				identToken("console"),
				dotToken(),
				identToken("log"),
				leftParenToken(),
				rightParenToken(),
			),
		},
		{
			name: "WithThreeDigitsDecimalParameter",
			code: Str("test(666)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("666"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitsDecimalParameter",
			code: Str("test(66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("66"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitDecimalParameter",
			code: Str("test(6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("6"),
				rightParenToken(),
			),
		},
		{
			name: "DecimalWithExponentParameter",
			code: Str("test(1e6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1e6"),
				rightParenToken(),
			),
		},
		{
			name: "DecimalWithUpperExponentParameter",
			code: Str("test(1E6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1E6"),
				rightParenToken(),
			),
		},
		{
			name: "WithSmallestRealDecimalParameter",
			code: Str("test(.1)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken(".1"),
				rightParenToken(),
			),
		},
		{
			name: "RealDecimalWithExponentParameter",
			code: Str("test(1.1e6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1.1e6"),
				rightParenToken(),
			),
		},
		{
			name: "RealDecimalWithUpperExponentParameter",
			code: Str("test(1.1E6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("1.1E6"),
				rightParenToken(),
			),
		},
		{
			name: "WithRealDecimalParameter",
			code: Str("test(6.6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				decimalToken("6.6"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitHexadecimalParameter",
			code: Str("test(0x6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0x6"),
				rightParenToken(),
			),
		},
		{
			name: "WithOneDigitUpperHexadecimalParameter",
			code: Str("test(0X6)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X6"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitHexadecimalParameter",
			code: Str("test(0x66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0x66"),
				rightParenToken(),
			),
		},
		{
			name: "WithTwoDigitUpperHexadecimalParameter",
			code: Str("test(0X66)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X66"),
				rightParenToken(),
			),
		},
		{
			name: "CommaSeparatedNumbersParameters",
			code: Str("test(0X6,0x7,0x78,0X69,8,69,669,6.9,.9,3e1,4E7,4e7)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				hexToken("0X6"),
				commaToken(),
				hexToken("0x7"),
				commaToken(),
				hexToken("0x78"),
				commaToken(),
				hexToken("0X69"),
				commaToken(),
				decimalToken("8"),
				commaToken(),
				decimalToken("69"),
				commaToken(),
				decimalToken("669"),
				commaToken(),
				decimalToken("6.9"),
				commaToken(),
				decimalToken(".9"),
				commaToken(),
				decimalToken("3e1"),
				commaToken(),
				decimalToken("4E7"),
				commaToken(),
				decimalToken("4e7"),
				rightParenToken(),
			),
		},
		{
			name: "CommaSeparatedNumbersAndStringsParameters",
			code: Str(`test("",5,"i",4,"k",6.6,0x5,"jssucks")`),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				stringToken(""),
				commaToken(),
				decimalToken("5"),
				commaToken(),
				stringToken("i"),
				commaToken(),
				decimalToken("4"),
				commaToken(),
				stringToken("k"),
				commaToken(),
				decimalToken("6.6"),
				commaToken(),
				hexToken("0x5"),
				commaToken(),
				stringToken("jssucks"),
				rightParenToken(),
			),
		},
		{
			name: "PassingIdentifierAsArg",
			code: Str("test(arg)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				identToken("arg"),
				rightParenToken(),
			),
		},
		{
			name: "PassingIdentifiersAsArg",
			code: Str("test(arg,arg2,i4k)"),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				identToken("arg"),
				commaToken(),
				identToken("arg2"),
				commaToken(),
				identToken("i4k"),
				rightParenToken(),
			),
		},
		{
			name: "CommaSeparatedEverything",
			code: Str(`test("",5,"i",4,"k",6.6,0x5,arg,"jssucks")`),
			want: tokens(
				identToken("test"),
				leftParenToken(),
				stringToken(""),
				commaToken(),
				decimalToken("5"),
				commaToken(),
				stringToken("i"),
				commaToken(),
				decimalToken("4"),
				commaToken(),
				stringToken("k"),
				commaToken(),
				decimalToken("6.6"),
				commaToken(),
				hexToken("0x5"),
				commaToken(),
				identToken("arg"),
				commaToken(),
				stringToken("jssucks"),
				rightParenToken(),
			),
		},
	})
}

func TestPosition(t *testing.T) {
	runTests(t, []TestCase{
		{
			name:          "MinusDecimal",
			code:          Str("-1"),
			checkPosition: true,
			want: tokens(minusTokenPos(1, 1), decimalTokenPos("1", 1, 2)),
		},
		{
			name:          "PlusDecimal",
			code:          Str("+1"),
			checkPosition: true,
			want: tokens(plusTokenPos(1, 1), decimalTokenPos("1", 1, 2)),
		},
		{
			name:          "PlusMinusDecimal",
			code:          Str("+-666"),
			checkPosition: true,
			want: tokens(plusTokenPos(1, 1), minusTokenPos(1, 2), decimalTokenPos("666", 1, 3)),
		},
	})
}

func TestIllegalIdentifiers(t *testing.T) {
	t.Skip("TODO")
}

func TestIllegalMemberAccess(t *testing.T) {

	runTests(t, []TestCase{
		{
			name: "CantAccessMemberThatStartsWithNumber",
			code: Str("test.123"),
			want: []lexer.Tokval{
				identToken("test"),
				dotToken(),
				illegalToken("123"),
			},
		},
		{
			name: "CantAccessMemberThatStartsWithDot",
			code: Str("test.."),
			want: []lexer.Tokval{
				identToken("test"),
				dotToken(),
				illegalToken("."),
			},
		},
	})
}

func TestIllegalNumericLiterals(t *testing.T) {

	corruptedHex := messStr(Str("0x01234"), 4)
	corruptedDecimal := messStr(Str("1234"), 3)
	corruptedNumber := messStr(Str("0"), 1)

	runTests(t, []TestCase{
		{
			name: "DecimalDuplicatedUpperExponentPart",
			code: Str("123E123E123"),
			want: []lexer.Tokval{
				illegalToken("123E123E123"),
			},
		},
		{
			name: "DecimalDuplicatedExponentPart",
			code: Str("123e123e123"),
			want: []lexer.Tokval{
				illegalToken("123e123e123"),
			},
		},
		{
			name: "RealDecimalDuplicatedUpperExponentPart",
			code: Str("123.1E123E123"),
			want: []lexer.Tokval{
				illegalToken("123.1E123E123"),
			},
		},
		{
			name: "RealDecimalDuplicatedExponentPart",
			code: Str("123.6e123e123"),
			want: []lexer.Tokval{
				illegalToken("123.6e123e123"),
			},
		},
		{
			name: "OnlyStartAsDecimal",
			code: Str("0LALALA"),
			want: []lexer.Tokval{
				illegalToken("0LALALA"),
			},
		},
		{
			name: "EndIsNotDecimal",
			code: Str("0123344546I4K"),
			want: []lexer.Tokval{
				illegalToken("0123344546I4K"),
			},
		},
		{
			name: "EmptyHexadecimal",
			code: Str("0x"),
			want: []lexer.Tokval{
				illegalToken("0x"),
			},
		},
		{
			name: "OnlyStartAsReal",
			code: Str("0.b"),
			want: []lexer.Tokval{
				illegalToken("0.b"),
			},
		},
		{
			name: "RealWithTwoDotsStartingWithDot",
			code: Str(".1.2"),
			want: []lexer.Tokval{
				illegalToken(".1.2"),
			},
		},
		{
			name: "RealWithTwoDots",
			code: Str("0.1.2"),
			want: []lexer.Tokval{
				illegalToken("0.1.2"),
			},
		},
		{
			name: "BifRealWithTwoDots",
			code: Str("1234.666.2342"),
			want: []lexer.Tokval{
				illegalToken("1234.666.2342"),
			},
		},
		{
			name: "EmptyHexadecimalUpperX",
			code: Str("0X"),
			want: []lexer.Tokval{
				illegalToken("0X"),
			},
		},
		{
			name: "LikeHexadecimal",
			code: Str("0b1234"),
			want: []lexer.Tokval{
				illegalToken("0b1234"),
			},
		},
		{
			name: "OnlyStartAsHexadecimal",
			code: Str("0xI4K"),
			want: []lexer.Tokval{
				illegalToken("0xI4K"),
			},
		},
		{
			name: "EndIsNotHexadecimal",
			code: Str("0x123456G"),
			want: []lexer.Tokval{
				illegalToken("0x123456G"),
			},
		},
		{
			name: "CorruptedHexadecimal",
			code: corruptedHex,
			want: []lexer.Tokval{
				illegalToken(corruptedHex.String()),
			},
		},
		{
			name: "CorruptedDecimal",
			code: corruptedDecimal,
			want: []lexer.Tokval{
				illegalToken(corruptedDecimal.String()),
			},
		},
		{
			name: "CorruptedNumber",
			code: corruptedNumber,
			want: []lexer.Tokval{
				illegalToken(corruptedNumber.String()),
			},
		},
	})
}

func TestNoOutputFor(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "EmptyString",
			code: Str(""),
			want: []lexer.Tokval{EOF},
		},
	})
}

func TestCorruptedInput(t *testing.T) {
	runTests(t, []TestCase{
		{
			name: "AtStart",
			code: messStr(Str(""), 0),
			want: []lexer.Tokval{illegalToken(messStr(Str(""), 0).String())},
		},
	})
}

func runTests(t *testing.T, testcases []TestCase) {

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tokensStream := lexer.Lex(tc.code)
			tokens := []lexer.Tokval{}

			for t := range tokensStream {
				tokens = append(tokens, t)
			}

			assertWantedTokens(t, tc, tokens)
		})
	}
}

func illegalToken(val string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Illegal,
		Value: Str(val),
	}
}

func assertWantedTokens(t *testing.T, tc TestCase, got []lexer.Tokval) {
	t.Helper()

	if len(tc.want) != len(got) {
		t.Errorf("wanted [%d] tokens, got [%d] tokens", len(tc.want), len(got))
		t.Fatalf("\nwant=%v\ngot= %v\nare not equal.", tc.want, got)
	}

	for i, w := range tc.want {
		g := got[i]
		if !w.Equal(g) {
			t.Errorf("\nwanted:\ntoken[%d][%v]\n\ngot:\ntoken[%d][%v]", i, w, i, g)
			t.Errorf("\nwanted:\n%v\ngot:\n%v\n", tc.want, got)
		}

		if tc.checkPosition {
			if !w.EqualPos(g) {
				t.Errorf("want=%+v\ngot=%+v\nare equal but dont have the same position", w, g)
			}
		}
	}
}

func messStr(s utf16.Str, pos uint) utf16.Str {
	// WHY: The go's utf16 package uses the replacement char everytime a some
	// encoding/decoding error happens, so we inject one on the uint16 array to simulate
	// encoding/decoding errors.
	// Not safe but the idea is to fuck up the string

	r := append(s[0:pos], uint16(unicode.ReplacementChar))
	r = append(r, s[pos:]...)
	return r
}

// prependOnTestCases will prepend the given tcase on each TestCase
// provided on tcases, generating a new array of TestCases.
//
// The array of TestCases is generated by prepending code and the
// wanted tokens from the given tcase on each test case on tcases.
// EOF should not be provided on the
// given tcase since it will be prepended on each test case inside given tcases.
func prependOnTestCases(tcase TestCase, tcases []TestCase) []TestCase {
	newcases := make([]TestCase, len(tcases))

	for i, t := range tcases {
		name := fmt.Sprintf("%s/%s", tcase.name, t.name)
		code := tcase.code.Append(t.code)
		want := append(tcase.want, t.want...)

		newcases[i] = TestCase{
			name: name,
			code: code,
			want: want,
		}
	}

	return newcases
}

func sfmt(format string, a ...interface{}) utf16.Str {
	return Str(fmt.Sprintf(format, a...))
}

func minusToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Minus,
		Value: Str("-"),
	}
}

func plusToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Plus,
		Value: Str("+"),
	}
}

func leftParenToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.LParen,
		Value: Str("("),
	}
}

func rightParenToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.RParen,
		Value: Str(")"),
	}
}

func minusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Minus,
		Value:  Str("-"),
		Line:   line,
		Column: column,
	}
}

func plusTokenPos(line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Plus,
		Value:  Str("+"),
		Line:   line,
		Column: column,
	}
}

func decimalTokenPos(dec string, line uint, column uint) lexer.Tokval {
	return lexer.Tokval{
		Type:   token.Decimal,
		Value:  Str(dec),
		Line:   line,
		Column: column,
	}
}

func decimalToken(dec string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Decimal,
		Value: Str(dec),
	}
}

func dotToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Dot,
		Value: Str("."),
	}
}

func hexToken(hex string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Hexadecimal,
		Value: Str(hex),
	}
}

func stringToken(s string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.String,
		Value: Str(s),
	}
}

func identToken(s string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Ident,
		Value: Str(s),
	}
}

func ltToken(s string) lexer.Tokval {
	return lexer.Tokval{
		Type:  token.LineTerminator,
		Value: Str(s),
	}
}

func commaToken() lexer.Tokval {
	return lexer.Tokval{
		Type:  token.Comma,
		Value: Str(","),
	}
}

func tokens(t ...lexer.Tokval) []lexer.Tokval {
	return append(t, EOF)
}