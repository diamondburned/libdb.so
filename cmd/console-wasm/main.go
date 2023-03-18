package main

import (
	"context"
	"io"

	"libdb.so/console"
	"libdb.so/console/programs"
)

var terminal *io.PipeReader // js reads from this
var input *io.PipeWriter    // js writes to this

var interp *console.Interpreter

func main() {
	rr, rw := io.Pipe()
	terminal = rr

	wr, ww := io.Pipe()
	input = ww

	io := console.IO{
		Stdin:  wr,
		Stdout: rw,
		Stderr: rw,
	}

	inst, err := console.NewInterpreter(io, console.InterpreterOpts{
		Programs: programs.All(),
	})
	if err != nil {
		panic(err)
	}

	interp = inst
	defer interp.Close()

	if err := interp.Run(context.Background()); err != nil {
		panic(err)
	}
}
