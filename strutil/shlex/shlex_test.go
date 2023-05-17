/*
Copyright 2012 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shlex

import (
	"errors"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var (
	// one two "three four" "five \"six\"" seven#eight # nine # ten
	// eleven 'twelve\'
	testString = "\\one two \"three four\" \"five \\\"six\\\"\" seven#eight # nine # ten\n eleven 'twelve\\' thirteen=13 fourteen/14"
)

func TestClassifier(t *testing.T) {
	classifier := newDefaultClassifier()
	tests := map[rune]runeTokenClass{
		' ':  spaceRuneClass,
		'"':  escapingQuoteRuneClass,
		'\'': nonEscapingQuoteRuneClass,
		'#':  commentRuneClass}
	for runeChar, want := range tests {
		got := classifier.ClassifyRune(runeChar)
		if got != want {
			t.Errorf("ClassifyRune(%v) -> %v. Want: %v", runeChar, got, want)
		}
	}
}

func TestTokenizer(t *testing.T) {
	testInput := strings.NewReader(testString)
	expectedTokens := []*Token{
		{WordToken, "one"},
		{WordToken, "two"},
		{WordToken, "three four"},
		{WordToken, "five \"six\""},
		{WordToken, "seven#eight"},
		{CommentToken, " nine # ten"},
		{WordToken, "eleven"},
		{WordToken, "twelve\\"},
		{WordToken, "thirteen=13"},
		{WordToken, "fourteen/14"}}

	tokenizer := NewTokenizer(testInput)
	for i, want := range expectedTokens {
		got, err := tokenizer.Next()
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf("Tokenizer.Next()[%v] of %q -> %v. Want: %v", i, testString, got, want)
		}
	}
}

func TestLexer(t *testing.T) {
	testInput := strings.NewReader(testString)
	expectedStrings := []string{"one", "two", "three four", "five \"six\"", "seven#eight", "eleven", "twelve\\", "thirteen=13", "fourteen/14"}

	lexer := NewLexer(testInput)
	for i, want := range expectedStrings {
		got, err := lexer.Next()
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("Lexer.Next()[%v] of %q -> %v. Want: %v", i, testString, got, want)
		}
	}
}

func TestSplit(t *testing.T) {
	want := []string{"one", "two", "three four", "five \"six\"", "seven#eight", "eleven", "twelve\\", "thirteen=13", "fourteen/14"}
	got, err := Split(testString)
	if err != nil {
		t.Error(err)
	}
	if len(want) != len(got) {
		t.Errorf("Split(%q) -> %v. Want: %v", testString, got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Split(%q)[%v] -> %v. Want: %v", testString, i, got[i], want[i])
		}
	}
}

func TestEOFAfterEscape(t *testing.T) {
	_, err := Split(testString + "\\")
	if err == nil {
		t.Error(err)
	}
}

func TestEOFInQuotingEscape(t *testing.T) {
	_, err := Split(`foo"`)
	if err == nil {
		t.Error(err)
	}

	_, err = Split(`foo'`)
	if err == nil {
		t.Error(err)
	}

	_, err = Split(`"foo\`)
	if err == nil {
		t.Error(err)
	}
}

func TestEOFInComment(t *testing.T) {
	got, err := Split("#")
	if err != nil {
		t.Error(err)
	}
	if len(got) > 1 {
		t.Errorf("Split(%q) -> %v", testString, got)
	}
}

type nastyReader struct{}

var errNastyReader = errors.New("foo")

func (*nastyReader) Read(_ []byte) (int, error) {
	return 0, errNastyReader
}

func TestNastyReader(t *testing.T) {
	l := NewLexer(&nastyReader{})
	_, err := l.Next()
	if err == nil {
		t.Errorf("expected an error, got nil instead")
	}
	if err != errNastyReader {
		t.Errorf("unexpected error")
	}
}

var quoteTests = []struct {
	name     string
	input    string
	expected string
}{{
	name:     `empty`,
	input:    ``,
	expected: `''`,
}, {
	name:     `quote safe`,
	input:    "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "@%_-+=:,./",
	expected: "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789" + "@%_-+=:,./",
}, {
	name:     `spaces`,
	input:    `foo bar xyz`,
	expected: `'foo bar xyz'`,
}, {
	name:     `double quotation`,
	input:    `foo"bar`,
	expected: `'foo"bar'`,
}, {
	name:     `backtick`,
	input:    "foo`bar",
	expected: "'foo`bar'",
}, {
	name:     `backslash`,
	input:    `foo\bar`,
	expected: `'foo\bar'`,
}, {
	name:     `unicode`,
	input:    "foo\xe9bar",
	expected: "'foo\xe9bar'",
}, {
	name:     `single quote with exclamation point`,
	input:    `foo!'bar'`,
	expected: `'foo!'"'"'bar'"'"''`,
}, {
	name:     `single quote with dollar`,
	input:    `'foo$'bar`,
	expected: `''"'"'foo$'"'"'bar'`,
}}

func TestQuote(t *testing.T) {
	for _, test := range quoteTests {
		t.Run(test.name, func(t *testing.T) {
			quoted := Quote(test.input)
			if quoted != test.expected {
				t.Errorf("expected %s, got %s", test.expected, quoted)
			}
		})
	}
}

var joinTests = []struct {
	name     string
	input    []string
	expected string
}{{
	name:     `space in first arg`,
	input:    []string{`a `, `b`},
	expected: `'a ' b`,
}, {
	name:     `space in last arg`,
	input:    []string{`a`, ` b`},
	expected: `a ' b'`,
}, {
	name:     `space as an arg`,
	input:    []string{`a`, ` `, `b`},
	expected: `a ' ' b`,
}, {
	name:     `empty arg`,
	input:    []string{`a`, ``, `b`},
	expected: `a '' b`,
}, {
	name:     `double quotes in arg`,
	input:    []string{`"a`, `b"`},
	expected: `'"a' 'b"'`,
}, {
	name:     `long args`,
	input:    []string{`x y`, `/foo/bar`},
	expected: `'x y' /foo/bar`,
}, {
	name:     `empty slice`,
	input:    []string{},
	expected: ``,
}, {
	name:     `nil slice`,
	input:    nil,
	expected: ``,
}}

func TestJoin(t *testing.T) {
	for _, test := range joinTests {
		t.Run(test.name, func(t *testing.T) {
			joined := Join(test.input)
			if joined != test.expected {
				t.Errorf("expected %s, got %s", test.expected, joined)
			}
		})
	}
}
