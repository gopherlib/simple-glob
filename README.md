# simple-glob

> Go Simple Globbing Library.
> 
> This repository was forked from https://github.com/gobwas/glob, but with some simplifications.

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

BenchmarkParseGlobGoogleURL
BenchmarkParseGlobGoogleURL-8         	  908487	      1281 ns/op
BenchmarkParseGlobAbc
BenchmarkParseGlobAbc-8               	 2062412	       582.0 ns/op
BenchmarkParseGlobDef
BenchmarkParseGlobDef-8               	 2255916	       532.2 ns/op
BenchmarkParseGlobAbdef
BenchmarkParseGlobAbdef-8             	13398868	        91.54 ns/op

BenchmarkGlobMatchGoogleURL_True
BenchmarkGlobMatchGoogleURL_True-8    	39037932	        30.07 ns/op
BenchmarkGlobMatchGoogleURL_False
BenchmarkGlobMatchGoogleURL_False-8   	80282545	        15.46 ns/op
BenchmarkGlobMatchAbc
BenchmarkGlobMatchAbc-8               	231060092	         5.215 ns/op
BenchmarkGlobMatchAbc_False
BenchmarkGlobMatchAbc_False-8         	351424537	         3.456 ns/op
BenchmarkGlobMatchDef_True
BenchmarkGlobMatchDef_True-8          	219895939	         5.467 ns/op
BenchmarkGlobMatchDef_Flase
BenchmarkGlobMatchDef_Flase-8         	349861196	         3.484 ns/op
BenchmarkGlobMatchAbdef_True
BenchmarkGlobMatchAbdef_True-8        	570146192	         2.117 ns/op
BenchmarkGlobMatchAbdef_Flase
BenchmarkGlobMatchAbdef_Flase-8       	569324104	         2.094 ns/op
```

## Syntax

Only `*` is supported as a wildcard to match any string other than a delimiter.
