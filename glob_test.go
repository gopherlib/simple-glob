package glob

import (
	"testing"
)

type test struct {
	pattern    string
	match      string
	should     bool
	delimiters []rune
}

func TestGlob(t *testing.T) {
	for _, test := range []test{
		{should: true, pattern: "* cat * eyes", match: "my cat has very bright eyes"},
		{should: true, pattern: "", match: ""},
		{should: true, pattern: `a*b`, match: "a*b"},
		{should: false, pattern: `a\*b`, match: "a*b"},
		{should: true, pattern: `a\*b`, match: `a\*b`},

		{should: true, pattern: "*ä", match: "åä"},
		{should: true, pattern: "abc", match: "abc"},
		{should: true, pattern: "a*c", match: "abc"},
		{should: true, pattern: "a*c", match: "a12345c"},
		{should: false, pattern: "a?c", match: "a1c"},
		{should: true, pattern: "a.*", match: "a.b.c"},
		{should: true, pattern: "a.b", match: "a.b", delimiters: []rune{'.'}},
		{should: true, pattern: "a.*", match: "a.b", delimiters: []rune{'.'}},
		{should: false, pattern: "a.**", match: "a.b.c", delimiters: []rune{'.'}},
		{should: false, pattern: "a.?.c", match: "a.b.c", delimiters: []rune{'.'}},
		{should: true, pattern: "a.*.c", match: "a.b.c", delimiters: []rune{'.'}},
		{should: false, pattern: "a.?.?", match: "a.b.c", delimiters: []rune{'.'}},
		{should: false, pattern: "a.*", match: "a.b.c", delimiters: []rune{'.'}},
		{should: false, pattern: "?at", match: "cat"},
		{should: true, pattern: "ca*", match: "ca*t"},
		{should: true, pattern: "ca*", match: "ca*"},
		{should: true, pattern: "c*t", match: "ca*aaaaaaaat"},
		{should: false, pattern: "?at", match: "fat"},
		{should: true, pattern: "*", match: "abc"},
		{should: false, pattern: `\*`, match: "*"},
		{should: false, pattern: "**", match: "a.b.c", delimiters: []rune{'.'}},

		{should: false, pattern: "?at", match: "at"},
		{should: false, pattern: "?at", match: "fat", delimiters: []rune{'f'}},
		{should: false, pattern: "*at", match: "fat", delimiters: []rune{'f'}},
		{should: false, pattern: "a.*", match: "a.b.c", delimiters: []rune{'.'}},
		{should: false, pattern: "a.?.c", match: "a.bb.c", delimiters: []rune{'.'}},
		{should: true, pattern: "a.*.c", match: "a.bb.c", delimiters: []rune{'.'}},
		{should: false, pattern: "*", match: "a.b.c", delimiters: []rune{'.'}},

		{should: true, pattern: "*test", match: "this is a test"},
		{should: true, pattern: "this*", match: "this is a test"},
		{should: true, pattern: "*is *", match: "this is a test"},
		{should: true, pattern: "*is*a*", match: "this is a test"},
		{should: true, pattern: "**test**", match: "this is a test"},
		{should: true, pattern: "**is**a***test*", match: "this is a test"},

		{should: false, pattern: "*is", match: "this is a test"},
		{should: false, pattern: "*no*", match: "this is a test"},
		{should: true, pattern: "*abc", match: "abcabc"},
		{should: true, pattern: "/*", match: "/rate"},

		{should: true, pattern: "*//*.example.com", match: "https://www.example.com"},
		{should: true, pattern: "*//*.example.com", match: "https://www.example.com", delimiters: []rune{'.'}},
		{should: true, pattern: "*//*example.com", match: "https://www.example.com"},
		{should: false, pattern: "*//*example.com", match: "https://www.example.com", delimiters: []rune{'.'}},
		{should: false, pattern: "*//*.example.com", match: "http://example.com"},
		{should: false, pattern: "*//*.example.com", match: "http://example.com.net"},
		{should: true, pattern: "*//*example.com", match: "http://example.com"},
	} {
		t.Run(test.pattern, func(t *testing.T) {
			g := MustCompile(test.pattern, test.delimiters...)
			result := g.Match(test.match)
			if result != test.should {
				t.Errorf(
					"pattern %q matching %q should be %v but got %v\n%s",
					test.pattern, test.match, test.should, result, g,
				)
			}
		})
	}
}

var (
	testPatterns = map[string]struct {
		pattern string
		text    string
	}{
		`google-true`:  {`https://*.google.*`, "https://account.google.com"},
		`google-false`: {`https://*.google.*`, "https://google.com"},
		`abc-true`:     {`abc*`, "abcdef"},
		`abc-false`:    {`abc*`, "af"},
		`def-true`:     {`*def`, "abcdef"},
		`def-false`:    {`*def`, "af"},
		`abef-true`:    {`ab*ef`, "abcdef"},
		`abef-false`:   {`ab*ef`, "af"},
	}
)

func BenchmarkParseGlobGoogleURL(b *testing.B) {
	pattern := testPatterns["google-true"]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MustCompile(pattern.pattern)
	}
}

func BenchmarkParseGlobAbc(b *testing.B) {
	pattern := testPatterns["abc-true"]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MustCompile(pattern.pattern)
	}
}

func BenchmarkParseGlobDef(b *testing.B) {
	pattern := testPatterns["def-true"]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MustCompile(pattern.pattern)
	}
}

func BenchmarkParseGlobAbdef(b *testing.B) {
	pattern := testPatterns["abdef-true"]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MustCompile(pattern.pattern)
	}
}

func BenchmarkGlobMatchGoogleURL_True(b *testing.B) {
	pattern := testPatterns["google-true"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchGoogleURL_False(b *testing.B) {
	pattern := testPatterns["google-false"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchAbc(b *testing.B) {
	pattern := testPatterns["abc-true"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchAbc_False(b *testing.B) {
	pattern := testPatterns["abc-false"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchDef_True(b *testing.B) {
	pattern := testPatterns["def-true"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchDef_Flase(b *testing.B) {
	pattern := testPatterns["def-false"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchAbdef_True(b *testing.B) {
	pattern := testPatterns["abdef-true"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}

func BenchmarkGlobMatchAbdef_Flase(b *testing.B) {
	pattern := testPatterns["abdef-false"]
	c := MustCompile(pattern.pattern)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Match(pattern.text)
	}
}
