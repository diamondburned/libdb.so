module libdb.so

go 1.20

replace mvdan.cc/sh/v3 => github.com/diamondburned/mvdan-sh/v3 v3.0.0-20230318131347-17d55f04e1ac

require (
	github.com/fatih/color v1.15.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.7.0
	mvdan.cc/sh/v3 v3.6.1-0.20230318112031-1e04c5bd318f
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
)
