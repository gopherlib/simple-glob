# simple-glob

> Go Simple Globbing Library.
> 
> This repository was forked from https://github.com/gobwas/glob, but with some simplifications.
> Only `*` is supported as a wildcard to match any string other than a delimiter.

## Install

```shell
    go get github.com/gopherlib/simple-glob
```

## Example

```go

package main

import "github.com/gopherlib/simple-glob"

func main() {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile("*.github.com")
	g.Match("api.github.com") // true

	// create new glob with set of delimiters as ["."]
	g = glob.MustCompile("api.*.com", '.')
	g.Match("api.github.com") // true
	g.Match("api.gi.hub.com") // false
}

```

## Performance

This library is created for compile-once patterns. This means, that compilation could take time, but
strings matching is done faster, than in case when always parsing template.

If you do not use compiled `glob.Glob` object, and do `g := glob.MustCompile(pattern); g.Match(...)` every time, then
your code will be much slower.

Run `go test -bench=.` from source root to see the benchmarks:

| Pattern              | Fixture                      | Match   | Speed (ns/op) |
|----------------------|------------------------------|---------|---------------|
| `https://*.google.*` | `https://account.google.com` | `true`  | 30.07         |
| `https://*.google.*` | `https://google.com`         | `false` | 15.46         |
| `abc*`               | `abcdef`                     | `true`  | 5.215         |
| `abc*`               | `af`                         | `false` | 3.456         |
| `*def`               | `abcdef`                     | `true`  | 5.467         |
| `*def`               | `af`                         | `false` | 3.484         |
| `ab*ef`              | `abcdef`                     | `true`  | 2.117         |
| `ab*ef`              | `af`                         | `false` | 2.094         |

```text
goos: darwin
goarch: arm64
pkg: github.com/gopherlib/simple-glob

BenchmarkParseGlobGoogleURL-8             904128              1278 ns/op            1400 B/op         39 allocs/op
BenchmarkParseGlobAbc-8                  2074610               580.4 ns/op           744 B/op         20 allocs/op
BenchmarkParseGlobDef-8                  2280055               528.0 ns/op           712 B/op         18 allocs/op
BenchmarkParseGlobAbdef-8               13565740                94.31 ns/op          256 B/op          3 allocs/op

BenchmarkGlobMatchGoogleURL_True-8      40695439                29.89 ns/op            0 B/op          0 allocs/op
BenchmarkGlobMatchGoogleURL_False-8     73579942                15.84 ns/op            0 B/op          0 allocs/op
BenchmarkGlobMatchAbc-8                 231631929                5.211 ns/op           0 B/op          0 allocs/op
BenchmarkGlobMatchAbc_False-8           348261666                3.472 ns/op           0 B/op          0 allocs/op
BenchmarkGlobMatchDef_True-8            219939248                5.459 ns/op           0 B/op          0 allocs/op
BenchmarkGlobMatchDef_Flase-8           350182290                3.419 ns/op           0 B/op          0 allocs/op
BenchmarkGlobMatchAbdef_True-8          571836799                2.093 ns/op           0 B/op          0 allocs/op
BenchmarkGlobMatchAbdef_Flase-8         573770670                2.095 ns/op           0 B/op          0 allocs/op
```

## Syntax

Only `*` is supported as a wildcard to match any string other than a delimiter.
