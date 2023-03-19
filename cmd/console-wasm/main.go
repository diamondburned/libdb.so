package main

import (
	"context"
	"io"
	"log"
	"os"
	"reflect"
	"runtime"
	"syscall/js"
	"unsafe"

	"libdb.so/cmd/internal/global"
	"libdb.so/console"
	"libdb.so/console/programs"
)

var input io.Writer // js writes to this

var startCh = make(chan struct{}, 1)
var terminal *console.Terminal

func main() {
	terminal = console.NewTerminal(newIO(), console.TerminalQuery{})

	ctx := context.Background()
	env := console.Environment{
		Terminal:   terminal,
		Programs:   programs.All(),
		Filesystem: global.Filesystem,
		Cwd:        global.InitialCwd,
	}

	interp, err := console.NewInterpreter(&env, console.InterpreterOpts{
		RunCommands: global.RC,
	})
	if err != nil {
		log.Panicln("cannot make new interpreter:", err)
	}

	global := js.Global()
	global.Set("console_write_stdin", js.FuncOf(write_stdin))
	global.Set("console_update_terminal", js.FuncOf(update_terminal))
	global.Set("console_start", js.FuncOf(start))

	<-startCh

	if err := interp.Run(ctx); err != nil {
		log.Panicln(err)
	}
}

func newIO() console.IO {
	wr, ww := io.Pipe()
	input = ww

	return console.IO{
		Stdin:  wr,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// start unblocks main and starts the interpreter loop. The JS side must have
// called update_terminal before calling this function.
func start(this js.Value, args []js.Value) any {
	startCh <- struct{}{}
	return nil
}

// write_stdin writes the given bytes pointer that is of len size into the stdin
// pipe. It returns false if the write failed. In most cases, the program should
// panic if that is the case.
//
// Note: the function MUST block until the write is complete. It also must not
// hold onto the bytes pointer after the write is complete.
func write_stdin(this js.Value, args []js.Value) any { // (string) => void
	s := args[0].String()
	b := unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
	input.Write(b)
	runtime.KeepAlive(s)
	return nil
}

func update_terminal(this js.Value, args []js.Value) any { // ({row, col, xpixel, ypixel, sixel}) => void
	terminal.UpdateQuery(console.TerminalQuery{
		Width:  args[0].Get("row").Int(),
		Height: args[0].Get("col").Int(),
		XPixel: args[0].Get("xpixel").Int(),
		YPixel: args[0].Get("ypixel").Int(),
		SIXEL:  args[0].Get("sixel").Bool(),
	})
	return nil
}
