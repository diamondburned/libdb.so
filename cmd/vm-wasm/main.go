package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"syscall/js"
	"unsafe"

	"libdb.so/cmd/internal/global"
	"libdb.so/internal/nsfw"
	"libdb.so/vm"
	"libdb.so/vm/programs"
	"libdb.so/vm/programs/neofetch"
	"libdb.so/vm/rwfs"
	"libdb.so/vm/rwfs/httpfs"
	"libdb.so/vm/rwfs/kvfs"
)

var gitrev string

func init() {
	if gitrev != "" {
		neofetch.OverrideGitRevision(gitrev)
	}
}

var input io.Writer // js writes to this
var startCh = make(chan struct{}, 1)
var terminal vm.Terminal
var publicFS *httpfs.FS

func main() {
	wr, ww := io.Pipe()
	input = ww

	vmIO := vm.IO{
		Stdin:  wr,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	terminal = vm.NewTerminal(vmIO, vm.TerminalQuery{})

	{
		global := js.Global()
		global.Set("vm_write_stdin", js.FuncOf(write_stdin))
		global.Set("vm_update_terminal", js.FuncOf(update_terminal))
		global.Set("vm_start", js.FuncOf(start))
		global.Set("vm_set_public_fs", js.FuncOf(set_public_fs))
	}

	<-startCh

	ctx := context.Background()
	env := vm.Environment{
		Terminal: terminal,
		Programs: programs.All(),
		Filesystem: rwfs.OverlayFS(
			kvfs.New(kvfs.LocalStorage()),
			rwfs.ReadOnlyFS(nsfw.WrapFS(publicFS)),
		),
		Cwd:     global.InitialCwd,
		Environ: global.InitialEnv,
	}

	interp, err := vm.NewInterpreter(&env, vm.InterpreterOpts{
		RunCommands: global.RC,
		Prompt:      global.PromptColored(),
		IgnoreEOF:   true,
	})
	if err != nil {
		log.Panicln("cannot make new interpreter:", err)
	}

	if err := interp.Run(ctx); err != nil {
		log.Panicln(err)
	}

	log.Println("interpreter exited. Bye!")
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

// set_public_fs sets the public file system to the given JSON string.
func set_public_fs(this js.Value, args []js.Value) any { // (string, string) => void
	jsonStr := args[0].String()
	basePath := args[1].String()

	var tree httpfs.FileTree
	if err := json.Unmarshal([]byte(jsonStr), &tree); err != nil {
		log.Panicln("cannot unmarshal public fs:", err)
	}

	publicFS = httpfs.New(*http.DefaultClient, tree, basePath)
	return nil
}

func update_terminal(this js.Value, args []js.Value) any { // ({row, col, xpixel, ypixel, sixel}) => void
	terminal.UpdateQuery(vm.TerminalQuery{
		Width:  args[0].Get("col").Int(),
		Height: args[0].Get("row").Int(),
		XPixel: args[0].Get("xpixel").Int(),
		YPixel: args[0].Get("ypixel").Int(),
		SIXEL:  args[0].Get("sixel").Bool(),
	})
	return nil
}
