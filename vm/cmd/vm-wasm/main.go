package main

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"syscall/js"
	"unsafe"

	"libdb.so/vm"
	"libdb.so/vm/fs/httpfs"
	"libdb.so/vm/fs/kvfs"
	"libdb.so/vm/fs/rwfs"
	"libdb.so/vm/global"
	"libdb.so/vm/internal/nsfw"
	"libdb.so/vm/programs"
	"libdb.so/vm/programs/neofetch"
)

var gitrev string

func init() {
	if gitrev != "" {
		neofetch.OverrideGitRevision(gitrev)
	}
}

var (
	startCh = make(chan struct{}, 1)
	stopCh  = make(chan struct{})

	input     io.Writer // js writes to this
	cancel    context.CancelFunc
	terminal  vm.Terminal
	extraFSes []fs.FS
)

func main() {
	wr, ww := io.Pipe()
	input = ww

	vmIO := vm.IO{
		Stdin:  wr,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	terminal = vm.NewTerminal(vmIO, vm.TerminalQuery{})

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	jsExport()

	<-startCh
	defer close(stopCh)

	env := vm.Environment{
		Terminal: terminal,
		Programs: programs.All(),
		Filesystem: rwfs.OverlayFS(
			kvfs.New(kvfs.LocalStorage()),
			global.RootFS(extraFSes...)...,
		),
		Cwd:     global.InitialCwd,
		Environ: global.InitialEnv,
	}

	interp, err := vm.NewInterpreter(&env, vm.InterpreterOpts{
		Prompt:      global.PromptColored(),
		IgnoreEOF:   true,
		RunCommands: ". .shellrc",
	})
	if err != nil {
		log.Panicln("cannot make new interpreter:", err)
	}

	if err := interp.Run(ctx); err != nil {
		log.Panicln(err)
	}

	log.Println("interpreter exited. Bye!")
}

func jsExport() {
	global := js.Global()
	global.Set("vm_write_stdin", js.FuncOf(write_stdin))
	global.Set("vm_update_terminal", js.FuncOf(update_terminal))
	global.Set("vm_set_public_fs", js.FuncOf(set_public_fs))
	global.Set("vm_add_public_fs", js.FuncOf(add_public_fs))
	global.Set("vm_add_public_fs_url", js.FuncOf(add_public_fs_url))
	global.Set("vm_stop", js.FuncOf(stop))
	global.Set("vm_start", js.FuncOf(start))
}

func jsUnexport() {
	global := js.Global()
	global.Set("vm_write_stdin", js.Undefined())
	global.Set("vm_update_terminal", js.Undefined())
	global.Set("vm_set_public_fs", js.Undefined())
	global.Set("vm_add_public_fs", js.Undefined())
	global.Set("vm_add_public_fs_url", js.Undefined())
	global.Set("vm_start", js.Undefined())
	global.Set("vm_stop", js.Undefined())
}

// start unblocks main and starts the interpreter loop. The JS side must have
// called update_terminal before calling this function.
func start(this js.Value, args []js.Value) any {
	startCh <- struct{}{}
	return nil
}

// stop stops the interpreter loop.
func stop(this js.Value, args []js.Value) any {
	cancel()
	<-stopCh
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
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	input.Write(b)
	runtime.KeepAlive(s)
	return nil
}

// set_public_fs sets the public file system to the given JSON string.
//
// Deprecated: use [add_public_fs] instead. This function acts the same as
// that function.
func set_public_fs(this js.Value, args []js.Value) any { // (string, string) => void
	return add_public_fs(this, args)
}

// add_public_fs adds the given httpfs JSON string to the public file system.
// It is added on top as a FS overlay.
func add_public_fs(this js.Value, args []js.Value) any { // (string, string) => void
	jsonStr := args[0].String()
	basePath := args[1].String()

	var tree httpfs.FileTree
	if err := json.Unmarshal([]byte(jsonStr), &tree); err != nil {
		log.Panicln("cannot unmarshal public fs:", err)
	}

	var fs fs.FS
	fs = httpfs.New(http.DefaultClient, httpfs.FileTreeRoot{Tree: tree}, basePath)
	fs = nsfw.WrapFS(fs)

	extraFSes = append(extraFSes, fs)
	return nil
}

// add_public_fs_url adds the given URL pointing to an fs.json to the public
// file system. It is added on top as a FS overlay.
func add_public_fs_url(this js.Value, args []js.Value) any {
	url := args[0].String()

	var fs fs.FS
	fs = httpfs.NewFromURL(http.DefaultClient, url)
	fs = nsfw.WrapFS(fs)

	extraFSes = append(extraFSes, fs)
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
