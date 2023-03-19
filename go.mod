module libdb.so

go 1.20

replace mvdan.cc/sh/v3 => github.com/diamondburned/mvdan-sh/v3 v3.0.0-20230318131347-17d55f04e1ac

//replace github.com/peterh/liner => github.com/diamondburned/peterh-liner v0.0.0-20230319000726-e4c9392f4efc
replace github.com/peterh/liner => ../liner

require (
	github.com/alecthomas/assert v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.15.0
	github.com/mattn/go-runewidth v0.0.3
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v3 v3.0.0-alpha2
	golang.org/x/crypto v0.7.0
	mvdan.cc/sh/v3 v3.6.1-0.20230318112031-1e04c5bd318f
)

require (
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.0.0-20210801044451-80ca428c5142 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
)
